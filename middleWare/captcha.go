package middleWare

import (
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

func GetCaptcha(c *gin.Context) {
	id := captcha.New()
	c.Header("Content-Type", "image/png")
	if captcha.WriteImage(c.Writer, id, captcha.StdWidth, captcha.StdHeight) != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "验证码绘制失败"})
		return
	}
	c.Header("X-Captcha-Id", id)
}

func Captcha(c *gin.Context) {
	id, userInput := c.PostForm("captchaID"), c.PostForm("captchaValue")
	if !captcha.VerifyString(id, userInput) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "验证码不正确"})
		return
	}
	c.Next()
}
