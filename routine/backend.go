package routine

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
	r.GET("/avatar", h.GetAvatar)
	r.GET("/getmail", h.Mailauth)
	r.GET("/getip", h.Captcha, h.AuthJwt, h.GetIP)

	//内容类型application/json的POST请求的路由逻辑
	jsonr := r.Group("/", conType("application/json"))
	jsonr.POST("/login", h.Captcha, h.Login)
	jsonr.POST("/signup", h.Signup)
	jsonr.POST("/rcon", h.Captcha, h.AuthJwt, h.Rcon)
	jsonr.POST("/gpt", h.RequestHeadersLogger(), h.AuthJwt, h.Gpt)
	jsonr.POST("/source", h.Captcha, h.AuthJwt, h.SetText)
	jsonr.POST("/gamebind", h.Captcha, h.AuthJwt, h.AuthBindCode)

	//启动后端
	r.Run(":" + co.Config.BackPort)
}

// 创建检查内容类型中间件的工厂函数
func conType(allowedType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Content-Type") != allowedType {
			util.Error(c, 400, "不接受此请求体格式，请使用"+allowedType, nil)
			return
		}
	}
}
