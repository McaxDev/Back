package routine

import (
	co "github.com/McaxDev/Back/config"
	h "github.com/McaxDev/Back/handler"
	"github.com/gin-gonic/gin"
)

func Backend() {
	r := gin.Default()
	r.Use(h.LogToSQL)
	r.GET("/captcha", h.GetCaptcha)
	r.GET("/challenge", h.GetChallenge)
	r.GET("/status", h.Status)
	r.GET("/prompt", h.Prompt)
	r.GET("/source", h.GetText)
	r.GET("/variable", h.Variable)
	r.POST("/login", h.Captcha, h.Login)
	r.POST("/signup", h.Captcha, h.Signup)
	r.GET("/getip", h.Captcha, h.Jwt, h.GetIP)
	r.POST("/rcon", h.Captcha, h.Jwt, h.Rcon)
	r.POST("/gpt", h.Captcha, h.Jwt, h.Gpt)
	r.POST("/source", h.Captcha, h.Jwt, h.SetText)
	r.Run(":" + co.Config.BackPort)
}
