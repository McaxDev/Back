package routine

import (
	conf "github.com/McaxDev/Back/config"
	hdlr "github.com/McaxDev/Back/handler"
	mid "github.com/McaxDev/Back/middleWare"
	"github.com/gin-gonic/gin"
)

func Backend() {
	r := gin.Default()
	r.GET("/captcha", mid.GetCaptcha)
	r.GET("/challenge", mid.GetChallenge)
	r.GET("/status", hdlr.Status)
	r.GET("/prompt", hdlr.Prompt)
	r.GET("/variable", hdlr.Variable)
	r.POST("/rcon", mid.Captcha, mid.Challenge(conf.Config.RconPwd), hdlr.Rcon)
	r.POST("/gpt", mid.Captcha, hdlr.Gpt)
	r.POST("/register", mid.Captcha, hdlr.Register)
	r.Run(":8080")
}
