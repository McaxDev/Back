package handler

import (
	"github.com/McaxDev/Back/util"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

func GetCaptcha(c *gin.Context) {
	id := captcha.New()
	c.Header("Content-Type", "image/png")
	c.Header("X-Captcha-Id", id)
	err := captcha.WriteImage(c.Writer, id, captcha.StdWidth, captcha.StdHeight)
	if err != nil {
		util.Error(c, 500, "验证码绘制失败", err)
		return
	}
}

func Captcha(c *gin.Context) {
	id, userInput := c.PostForm("captchaID"), c.PostForm("captchaValue")
	if !captcha.VerifyString(id, userInput) {
		util.Error(c, 400, "验证码不正确", nil)
		return
	}
	c.Next()
}
