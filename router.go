package main

import (
	"time"

	co "github.com/McaxDev/Back/config"
	h "github.com/McaxDev/Back/handler"
	"github.com/McaxDev/Back/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Backend() {

	//创建一个路由器
	r := gin.Default()

	//允许CORS跨域
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//全局应用记录日志到数据库的中间件
	r.Use(h.LogToSQL)

	//GET请求的路由逻辑
	r.GET("/captcha", h.GetCaptcha)
	r.GET("/challenge", h.GetChallenge)
	r.GET("/status", h.Status)
	r.GET("/prompt", h.Prompt)
	r.GET("/source", h.GetText)
	r.GET("/srvinfo", h.SrvInfo)
	r.GET("/getmail", h.RateLimit(60), h.Mailauth)

	// 对POST请求检查内容类型application/json的路由逻辑
	jsonr := r.Group("/", conType("application/json"))
	jsonr.POST("/login", h.Captcha, h.Login)
	jsonr.POST("/signup", h.Signup)

	// 内容为json而且要求jwt的请求的路由逻辑
	jsonjwtr := jsonr.Group("/", h.AuthJwt)
	jsonjwtr.GET("/autologin", h.AutoLogin)
	jsonjwtr.GET("/coin", h.Coin)
	jsonjwtr.GET("/gptutil", h.GptUtil)
	jsonjwtr.GET("/playerdata", h.PlayerData)
	jsonjwtr.GET("/getsms", h.RateLimit(60), h.SMS)
	jsonjwtr.GET("/change/username", h.ChangeUsername)
	jsonjwtr.GET("/change/gamename", h.ChangeGamename)
	jsonjwtr.POST("/change/password", h.ChangePwd)
	jsonjwtr.POST("/change/email", h.ChangeMail)
	jsonjwtr.POST("/gpt", h.Gpt)

	// 内容为json而且要求jwt而且要求人机验证的路由逻辑
	jsonjwtcapr := jsonjwtr.Group("/", h.Captcha)
	jsonjwtcapr.GET("/getip", h.GetIP)
	jsonjwtcapr.POST("/rcon", h.Rcon)
	jsonjwtcapr.POST("/source", h.SetText)
	jsonjwtcapr.POST("/gamebind", h.AuthBindCode)

	//启动后端
	r.RunTLS(":"+co.Config.BackPort, co.Config.SSL["pem"], co.Config.SSL["key"])
}

// 创建检查内容类型中间件的工厂函数
func conType(allowed string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" && c.Request.Header.Get("Content-Type") != allowed {
			util.Error(c, 400, "不接受此请求体格式，请使用"+allowed, nil)
			return
		}
	}
}
