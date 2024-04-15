package handler

import (
	"log"

	"github.com/gin-gonic/gin"
)

// RequestHeadersLogger 中间件打印请求中的所有头信息
func RequestHeadersLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Headers for %s %s:", c.Request.Method, c.Request.URL.Path)
		for k, v := range c.Request.Header {
			log.Printf("%s: %v", k, v)
		}
		c.Next()
	}
}
