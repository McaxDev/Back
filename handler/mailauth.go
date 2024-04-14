package handler

import (
	"fmt"
	"net/smtp"
	"sync"
	"time"

	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 记录已发送的邮箱的map
var Mailsent = make(map[string]MailStruct)
var IpTimeMap = make(map[string]time.Time)
var mu sync.Mutex

// 定义上面map的值的结构体，接受者和过期时间
type MailStruct struct {
	Receiver string
	Expiry   time.Time
}

func Mailauth(c *gin.Context) {

	//从查询字符串参数获得用户邮箱
	receiver := c.Query("receiver")

	//检查同一个IP是否在一分钟内重复请求
	clientip := c.ClientIP()
	mu.Lock()
	if time.Now().Before(IpTimeMap[clientip]) {
		lefttime := IpTimeMap[clientip].Sub(time.Now()).Seconds()
		mu.Unlock()
		util.Error(c, 400, fmt.Sprintf("请求频繁，请%.0f秒后重试", lefttime), nil)
		return
	}
	IpTimeMap[clientip] = time.Now().Add(time.Minute)
	mu.Unlock()

	//检查用户输入的邮箱是否已经被注册
	if err := co.DB.First(&co.User{}, "email = ?", receiver).Error; err == nil {
		util.Error(c, 400, "这个邮箱已经注册过了", err)
		return
	}

	//生成六位数验证码字符串
	authcode, err := util.RandStr(6)
	if err != nil {
		util.Error(c, 500, "服务器生成验证码失败", err)
		return
	}

	//生成过期时间并将验证码存到map里面
	expiry := time.Now().Add(10 * time.Minute)
	Mailsent[authcode] = MailStruct{
		Receiver: receiver,
		Expiry:   expiry,
	}

	//向请求者发送邮件
	fmttedExp := expiry.Format("2006-01-02 15:04")
	conf := co.Config.SMTPConfig
	dest := conf.Srv + ":" + conf.Port
	auth := smtp.PlainAuth("", conf.Mail, conf.Pwd, conf.Srv)
	to := []string{receiver}
	content := mailContent(receiver, authcode, fmttedExp)
	if err = smtp.SendMail(dest, auth, conf.Mail, to, content); err != nil {
		util.Error(c, 500, "邮件发送失败", err)
		return
	}
	util.Info(c, 200, "邮件发送成功", nil)
}

// 生成邮件界面标记语言
func mailContent(receiver, authcode, exp string) []byte {
	return []byte(
		"To: " + receiver + "\r\n" +
			"Subject: 验证码邮件\r\n" +
			"MIME-Version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
			fmt.Sprintf(`
<html>
<head>
  <style>
    body { font-family: 'Arial', sans-serif; background-color: #f4f4f4; margin: 0; padding: 20px; }
    .email-container { background-color: #ffffff; padding: 20px; margin: auto; max-width: 600px; border-radius: 8px; box-shadow: 0 0 10px rgba(0,0,0,0.1); }
    .header { font-size: 20px; margin-bottom: 20px; }
    .content { font-size: 16px; color: #333333; }
    .footer { font-size: 12px; color: #777777; margin-top: 20px; }
  </style>
</head>
<body>
  <div class="email-container">
    <p class="header">验证码</p>
    <p class="content">尊敬的用户：</p>
    <p class="content">您的验证码是<strong>%s</strong>，有效期至%s。</p>
    <p class="content">请在有效期内使用验证码进行验证。</p>
    <p class="footer">此邮件由系统自动发送，请勿直接回复。</p>
  </div>
</body>
</html>
`, authcode, exp))
}

// 验证邮箱的函数
func AuthMail(authcode, receiver string) bool {
	temp := Mailsent[authcode]
	if temp.Receiver != receiver || time.Now().After(temp.Expiry) {
		return false
	}
	delete(Mailsent, authcode)
	return true
}

// 清理过期的已发送邮件
func ClearExpiredMailsent() {
	now := time.Now()
	for key, mail := range Mailsent {
		if now.After(mail.Expiry) {
			delete(Mailsent, key)
		}
	}
}
