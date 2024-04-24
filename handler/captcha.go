package handler

import (
	"github.com/McaxDev/Back/util"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

// 获取captcha验证码
func GetCaptcha(c *gin.Context) {

	// 生成验证码ID
	id := captcha.New()

	// 设置响应头
	c.Header("Content-Type", "image/png")
	c.Header("X-Captcha-Id", id)

	// 发送响应体图片
	if err := captcha.WriteImage(
		c.Writer, id, captcha.StdWidth, captcha.StdHeight,
	); err != nil {
		util.Error(c, 500, "验证码绘制失败", err)
		return
	}
}

// 验证captcha验证码
func Captcha(c *gin.Context) {

	// 将请求体绑定到结构体
	var req struct {
		CaptchaID    string `json:"captchaID"`
		CaptchaValue string `json:"captchaValue"`
	}
	if err := util.BindReq(c, &req); err != nil {
		util.Error(c, 400, "无法将你的请求体绑定到结构体", err)
		return
	}

	//验证captcha验证码
	if !captcha.VerifyString(req.CaptchaID, req.CaptchaValue) {
		util.Error(c, 400, "验证码不正确", nil)
		return
	}
}
