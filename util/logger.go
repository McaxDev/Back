package util

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 向客户端发送错误
func Error(c *gin.Context, status int, msg string, err error) {
	if err != nil {
		c.Set("error", err)
	}
	c.AbortWithStatusJSON(status, Resp(msg, nil))
}

// 向客户端发送正常信息
func Info(c *gin.Context, status int, msg string, data any) {
	c.AbortWithStatusJSON(status, Resp(msg, data))
}

// 构造响应格式的映射
func Resp(msg string, data any) gin.H {
	return gin.H{"msg": msg, "data": data}
}

func ErrToStr(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func DbQueryError(c *gin.Context, err error, message string) {
	if err == gorm.ErrRecordNotFound {
		Error(c, 400, message, err)
	} else {
		Error(c, 500, "查询失败", err)
	}
}
