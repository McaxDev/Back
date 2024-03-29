package util

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/entity"
	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, status int, msg string, err error) {
	logAndPrint("ERROR", "\033[31m", status, msg, nil)
	c.AbortWithStatusJSON(status, gin.H{"msg": msg, "data": nil})
}

func Info(c *gin.Context, status int, msg string, data map[string]interface{}) {
	logAndPrint("INFO", "\033[32m", status, msg, nil)
	c.JSON(status, gin.H{"msg": msg, "data": data})
}

func Warn(c *gin.Context, status int, msg string, err error) {
	logAndPrint("WARN", "\033[33m", status, msg, err)
	c.AbortWithStatusJSON(status, gin.H{"msg": msg, "data": nil})
}

func Fatal(msg string, err error) {
	logAndPrint("FATAL", "\033[35m", 0, msg, err)
	os.Exit(1)
}

func logAndPrint(level string, color string, status int, msg string, err error) {
	errString := ""
	if err != nil {
		errString = err.Error()
	}
	DBlog := entity.Log{
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Level:   level,
		Status:  status,
		Message: msg,
		Error:   errString,
	}
	if dbErr := config.DB.Create(&DBlog).Error; dbErr != nil {
		log.Println("将日志存储到数据库失败：" + dbErr.Error())
	}
	colored := color + level + " " + strconv.Itoa(status) + "\033[0m"
	log.Println(fmt.Sprintf("%s %s %s", colored, msg, errString))
}
