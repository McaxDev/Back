package util

import (
	"log"
	"time"

	co "github.com/McaxDev/Back/config"
	"github.com/gin-gonic/gin"
)

func Res(msg string, data gin.H) gin.H {
	return gin.H{"msg": msg, "data": data}
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

func Data(pairs ...interface{}) gin.H {
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
