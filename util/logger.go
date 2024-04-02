package util

import (
	"log"
	"time"

	co "github.com/McaxDev/Back/config"
	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, status int, msg string, err error) {
	c.AbortWithStatusJSON(status, gin.H{"msg": msg, "data": nil})
}

func Info(c *gin.Context, status int, msg string, data map[string]interface{}) {
	c.JSON(status, gin.H{"msg": msg, "data": data})
}

func LogToSQL(c *gin.Context, level string, duration time.Duration) {
	errString := ""
	if UnknownErr, exist := c.Get("error"); exist {
		if err, ok := UnknownErr.(error); ok {
			errString = errToStr(err)
		}
	}
	DBlog := co.Log{
		Status:   c.Writer.Status(),
		Error:    errString,
		Method:   c.Request.Method,
		Path:     c.Request.URL.Path,
		Source:   c.ClientIP(),
		Duration: duration,
	}
	if dbErr := co.DB.Create(&DBlog).Error; dbErr != nil {
		log.Println("将日志存储到数据库失败：" + dbErr.Error())
	}
}

func MyMap(pairs ...interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	n := len(pairs)
	if n%2 != 0 {
		n -= 1
	}
	for i := 0; i < n; i += 2 {
		if key, ok := pairs[i].(string); ok {
			result[key] = pairs[i+1]
		}
	}
	return result
}

func errToStr(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
