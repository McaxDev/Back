package handler

import (
	"fmt"
	"net/smtp"
	"sync"
	"time"

	"github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 记录已发送的邮箱的map
var mailSent = make(map[string]MailStruct)
var ipTimeMap = make(map[string]time.Time)
var mu sync.Mutex

// 定义上面map的值的结构体，接受者和过期时间
type MailStruct struct {
	Receiver string
	Expiry   time.Time
}

func Mailauth(c *gin.Context) {

	//从查询字符串参数获得用户邮箱
	receiver := c.Query("receiver")
	to := []string{receiver}

	//检查同一个IP是否在一分钟内重复请求
	clientip := c.ClientIP()
	mu.Lock()
	if time.Now().Before(ipTimeMap[clientip].Add(time.Minute)) {
		lefttime := ipTimeMap[clientip].Add(time.Minute).Sub(time.Now()).Seconds()
		mu.Unlock()
		util.Error(c, 400, fmt.Sprintf("请求频繁，请%.0f秒后重试", lefttime), nil)
		return
	}
	ipTimeMap[clientip] = time.Now()
	mu.Unlock()

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
	fmttedExp := expiry.Format("2006-01-02 15:04")
	message := []byte("To: " + receiver + "\r\n" +
		"Subject: 验证码邮件\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n" +
		fmt.Sprintf("你的验证码是%s，%v内有效。", authcode, fmttedExp))
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
