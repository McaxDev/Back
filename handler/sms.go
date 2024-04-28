package handler

import (
	"time"

	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	unisms "github.com/apistd/uni-go-sdk/sms"
	"github.com/gin-gonic/gin"
)

// 存储已发送的邮件的数据结构
var smsSent = make(map[string]smsStru)

type smsStru struct {
	AuthCode string
	Expiry   time.Time
}

// 短信验证码接口
func SMS(c *gin.Context) {

	// 从中间件JWT里获取用户身份
	user, err := BindJwt(c, "Balance")
	if err != nil {
		util.Error(c, 400, "无法读取你的用户信息", err)
		return
	}
	// 从查询字符串参数里获取收件人号码
	receiver := c.Query("telephone")

	// 检查这个手机号是否已经被注册过了
	if err := co.DB.First(&co.User{}, "telephone = ?", receiver).Error; err != nil {
		util.Error(c, 400, "这个手机号已经被注册过了", err)
		return
	}

	// 对用户进行扣费
	if err := user.Transact(-10); err != nil {
		util.Error(c, 400, "你钱不够，10币一次", err)
		return
	}

	// 生成验证码和短信内容
	authCode := util.RandStr(6)

	// 向请求者的手机号发送短信
	conf := co.Config
	client := unisms.NewClient(conf.SMS["ID"], conf.SMS["Secret"])

	// 构造短信内容
	message := unisms.BuildMessage()
	message.SetTo(receiver)
	message.SetSignature("Axolotland")
	message.SetTemplateId("pub_verif_ttl3")
	message.SetTemplateData(map[string]string{
		"code": authCode,
		"ttl":  "10",
	})

	// 发送短信
	res, err := client.Send(message)
	if err != nil {
		util.Error(c, 500, "短信发送失败", err)
		return
	}

	// 将验证码和发件人存储到内存里
	smsSent[receiver] = smsStru{
		AuthCode: authCode,
		Expiry:   time.Now().Add(10 * time.Minute),
	}

	// 返回成功响应
	util.Info(c, 200, "操作成功，请在十分钟内处理短信", res)
}
