package routine

import (
	hdlr "github.com/McaxDev/Back/handler"
	"github.com/gin-gonic/gin"
)

func Backend() {
	r := gin.Default()
	r.GET("/captcha", hdlr.GetCaptcha)
	r.GET("/challenge", hdlr.GetChallenge)
	r.GET("/status", hdlr.Status)
	r.GET("/prompt", hdlr.Prompt)
	r.GET("/variable", hdlr.Variable)
	r.POST("/rcon", hdlr.Captcha, hdlr.Rcon)
	r.POST("/gpt", hdlr.Captcha, hdlr.Gpt)
	r.POST("/register", hdlr.Captcha, hdlr.Register)
	r.Run(":8080")
}
