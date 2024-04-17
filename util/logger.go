package util

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Error(c *gin.Context, status int, msg string, err error) {
	c.Set("error", err)
	c.AbortWithStatusJSON(status, gin.H{"msg": msg, "data": nil})
}

func Info(c *gin.Context, status int, msg string, data any) {
	c.AbortWithStatusJSON(status, gin.H{"msg": msg, "data": data})
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
