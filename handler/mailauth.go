package handler

import (
	"fmt"
	"net/smtp"
	"time"

	"github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

var mailSent = make(map[string]MailStruct)

type MailStruct struct {
	Receiver string
	Expiry   time.Time
}

func Mailauth(c *gin.Context) {
	receiver := c.Query("receiver")
	authcode, err := util.RandStr(6)
	if err != nil {
		c.AbortWithStatusJSON(500, util.Res("服务器生成验证码失败", nil))
		return
	}
	expiry := time.Now().Add(10 * time.Minute)
	mailSent[authcode] = MailStruct{
		Receiver: receiver,
		Expiry:   expiry,
	}
	message := []byte(fmt.Sprintf("你的验证码是%s，%v内有效。", authcode, expiry))
	conf := config.Config.SMTPConfig
	auth := smtp.PlainAuth("", conf.Mail, conf.Pwd, conf.Srv)
	to := []string{receiver}
	err = smtp.SendMail(conf.Srv+":"+conf.Port, auth, conf.Mail, to, message)
	if err != nil {
		c.AbortWithStatusJSON(400, util.Res("邮件发送失败", nil))
		return
	}
	c.AbortWithStatusJSON(200, util.Res("邮件发送成功", nil))
}

func AuthMail(authcode, receiver string) bool {
	temp := mailSent[authcode]
	if temp.Receiver == receiver && time.Now().After(temp.Expiry) {
		return false
	}
	delete(mailSent, authcode)
	return true
}

func ClearexpiredMailSent() {
	now := time.Now()
	for key, mail := range mailSent {
		if now.After(mail.Expiry) {
			delete(mailSent, key)
		}
	}
}
