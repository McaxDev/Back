package util

import "github.com/gin-gonic/gin"

func Json(msg string, data interface{}) gin.H {
	return gin.H{
		"msg":  msg,
		"data": data,
	}
}
