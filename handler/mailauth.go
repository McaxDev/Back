package handler

import (
	"fmt"
	"net/smtp"
	"time"

	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

var mailSent = make(map[string]MailMemory)

type MailMemory struct {
	Authcode string
	Expiry   time.Time
	LastSent time.Time // 添加最后发送时间字段
}

func Mailauth(c *gin.Context) {
	mail := c.Query("mail")

	// 检查是否存在记录，并且距离上次发送时间小于1分钟
	if mm, found := mailSent[mail]; found && time.Since(mm.LastSent) < time.Minute {
		c.AbortWithStatusJSON(429, util.Res("request too frequent", nil))
		return
	}

	authcode, err := util.RandStr(6)
	if err != nil {
		c.AbortWithStatusJSON(500, util.Res("fail to generate: "+err.Error(), nil))
		return
	}

	expiry := time.Now().Add(10 * time.Minute)
	mailSent[mail] = MailMemory{
		Authcode: authcode,
		Expiry:   expiry,
		LastSent: time.Now(),
	}

	message := []byte(fmt.Sprintf("Authcode is %s, available in %v.", authcode, expiry))
	conf := co.Config.SMTPConfig
	auth := smtp.PlainAuth("", conf.SMTPmail, conf.SMTPpwd, conf.SMTPsrv)
	to := []string{mail}
	err = smtp.SendMail(conf.SMTPsrv+":"+conf.SMTPport, auth, conf.SMTPmail, to, message)
	if err != nil {
		c.AbortWithStatusJSON(400, util.Res("fail to send: "+err.Error(), nil))
		return
	}

	c.AbortWithStatusJSON(200, util.Res("send successfully", nil))
}

func ClearexpiredMailSent() {
	now := time.Now()
	for key, mail := range mailSent {
		if now.After(mail.Expiry) {
			delete(mailSent, key)
		}
	}
}
