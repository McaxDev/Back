package handler

import (
	"fmt"
	"net/smtp"
	"time"

	"github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 记录已发送的邮箱的map
var mailSent = make(map[string]MailStruct)

// 定义上面map的值的结构体，接受者和过期时间
type MailStruct struct {
	Receiver string
	Expiry   time.Time
}

func Mailauth(c *gin.Context) {

	//从查询字符串参数获得用户邮箱
	receiver := c.Query("receiver")
	to := []string{receiver}

	//生成六位数验证码字符串
	authcode, err := util.RandStr(6)
	if err != nil {
		util.Error(c, 500, "服务器生成验证码失败", err)
		return
	}

	//生成过期时间并将验证码存到map里面
	expiry := time.Now().Add(10 * time.Minute)
	mailSent[authcode] = MailStruct{
		Receiver: receiver,
		Expiry:   expiry,
	}

	//向请求者发送邮件
	message := []byte(fmt.Sprintf("你的验证码是%s，%v内有效。", authcode, expiry))
	conf := config.Config.SMTPConfig
	auth := smtp.PlainAuth("", conf.Mail, conf.Pwd, conf.Srv)
	if err = smtp.SendMail(conf.Srv+":"+conf.Port, auth, conf.Mail, to, message); err != nil {
		util.Error(c, 500, "邮件发送失败", err)
		return
	}
	util.Info(c, 200, "邮件发送成功", nil)
}

// 验证邮箱的函数
func AuthMail(authcode, receiver string) bool {
	temp := mailSent[authcode]
	if temp.Receiver != receiver || time.Now().After(temp.Expiry) {
		return false
	}
	delete(mailSent, authcode)
	return true
}

// 清理过期的已发送邮件
func ClearexpiredMailSent() {
	now := time.Now()
	for key, mail := range mailSent {
		if now.After(mail.Expiry) {
			delete(mailSent, key)
		}
	}
}
