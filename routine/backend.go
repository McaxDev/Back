package routine

import (
	hdlr "github.com/McaxDev/Back/handler"
	"github.com/gin-gonic/gin"
)

func Backend() {
	r := gin.Default()
	r.Use(hdlr.LogToSQL)
	r.GET("/captcha", hdlr.GetCaptcha)
	r.GET("/challenge", hdlr.GetChallenge)
	r.GET("/status", hdlr.Status)
	r.GET("/prompt", hdlr.Prompt)
	r.GET("/variable", hdlr.Variable)
	r.POST("/rcon", hdlr.Captcha, hdlr.Jwt, hdlr.Rcon)
	r.POST("/gpt", hdlr.Captcha, hdlr.Jwt, hdlr.Gpt)
	r.POST("/login", hdlr.Captcha, hdlr.Login)
	r.POST("/register", hdlr.Captcha, hdlr.Register)
	r.POST("/source", hdlr.Source)
	r.Run(":8080")
}
