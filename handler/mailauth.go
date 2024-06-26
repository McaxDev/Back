package handler

import (
	"bytes"
	"html/template"
	"net/smtp"
	"time"

	"github.com/McaxDev/Back/assets"
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 记录已发送的邮箱的map
var Mailsent = make(map[string]MailStruct)
var tmpl *template.Template

// 定义上面map的值的结构体，接受者和过期时间
type MailStruct struct {
	Receiver string
	Expiry   time.Time
}

func Mailauth(c *gin.Context) {

	//从查询字符串参数获得用户邮箱
	receiver := c.Query("receiver")

	//检查用户输入的邮箱是否已经被注册
	if err := co.DB.First(&co.User{}, "email = ?", receiver).Error; err == nil {
		util.Error(c, 400, "这个邮箱已经注册过了", err)
		return
	}

	//生成六位数验证码字符串
	authcode := util.RandStr(6)

	//生成过期时间并将验证码存到map里面
	expiry := time.Now().Add(10 * time.Minute)
	Mailsent[authcode] = MailStruct{
		Receiver: receiver,
		Expiry:   expiry,
	}

	// 获得客户端住址
	address, err := util.Locateip(c.ClientIP())
	if err != nil {
		util.Error(c, 500, "无法将解析你的地址", err)
		return
	}

	//向请求者发送邮件
	fmttedExp := expiry.Format("2006-01-02 15:04")
	conf := co.Config.SMTPConfig
	dest := conf["server"] + ":" + conf["port"]
	auth := smtp.PlainAuth("", conf["mail"], conf["password"], conf["server"])
	to := []string{receiver}
	content, err := mailContent(receiver, authcode, fmttedExp, address)
	if err != nil {
		util.Error(c, 500, "邮件内容创建失败", err)
		return
	}
	if err = smtp.SendMail(dest, auth, conf["mail"], to, content); err != nil {
		util.Error(c, 500, "邮件发送失败", err)
		return
	}
	util.Info(c, 200, "邮件发送成功", nil)
}

// 生成邮件内容
func mailContent(receiver, authcode, exp, address string) ([]byte, error) {

	// 从嵌入的文件系统加载和解析邮件模板
	tmpl, err := template.New("mail_template.html").ParseFS(assets.Fs, "mail_template.html")
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	data := struct {
		Receiver   string
		AuthCode   string
		Expiration string
		Location   string
	}{
		Receiver:   receiver,
		AuthCode:   authcode,
		Expiration: exp,
		Location:   address,
	}

	// 定义邮件头部信息并直接创建为字节切片
	headers := []byte(
		"From: Axolotland Gaming Club <axolotland@163.com>\r\n" +
			"To: " + receiver + "\r\n" +
			"Subject: 验证码邮件\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"\r\n",
	)

	// 直接写入头部信息
	buf.Write(headers)

	// 执行模板并将生成的HTML添加到邮件内容
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
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
