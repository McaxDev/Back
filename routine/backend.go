package routine

import (
	"time"

	co "github.com/McaxDev/Back/config"
	h "github.com/McaxDev/Back/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Backend() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // 或者指定允许的域名
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(h.LogToSQL)
	r.GET("/captcha", h.GetCaptcha)
	r.GET("/challenge", h.GetChallenge)
	r.GET("/status", h.Status)
	r.GET("/prompt", h.Prompt)
	r.GET("/source", h.GetText)
	r.GET("/variable", h.Variable)
	r.GET("/avatar", h.GetAvatar)
	r.GET("/getmail", h.Mailauth)
	r.POST("/login", h.Login)
	r.POST("/signup", h.Captcha, h.Signup)
	r.GET("/getip", h.Captcha, h.AuthJwt, h.GetIP)
	r.POST("/rcon", h.Captcha, h.AuthJwt, h.Rcon)
	r.POST("/gpt", h.AuthJwt, h.AskGpt)
	r.POST("/source", h.Captcha, h.AuthJwt, h.SetText)
	r.POST("/gamebind", h.Captcha, h.AuthJwt, h.AuthBindCode)
	r.Run(":" + co.Config.BackPort)
}
