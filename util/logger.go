package util

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/McaxDev/Back/config"
	"github.com/gin-gonic/gin"
)

type ErrorLog struct {
	ID      int    `gorm:"primary_key"`
	Time    string `gorm:"type:datatime"`
	Level   string `gorm:"type:varchar(50)"`
	Status  int    `gorm:"type:int"`
	Message string `gorm:"type:text"`
}

func Error(c *gin.Context, status int, msg string, err error) {
	c.AbortWithStatusJSON(status, gin.H{"msg": msg, "data": nil})
	errorMsg := strconv.Itoa(status) + " " + msg
	if err != nil {
		errorMsg += ": " + err.Error()
	}
	log.Println(errorMsg)
	errorLog := ErrorLog{
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Level:   "ERROR",
		Status:  status,
		Message: errorMsg,
	}
	config.DB.Create(&errorLog)
}

func Info(c *gin.Context, status int, msg string, data map[string]interface{}) {
	ginData := gin.H{"msg": msg, "data": data}
	c.JSON(status, ginData)
	msg = strconv.Itoa(status) + " " + msg
	if jsonData, err := json.Marshal(ginData); err != nil {
		msg += " " + string(jsonData)
	}
	log.Println(msg)
	errorLog := ErrorLog{
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Level:   "INFO",
		Status:  status,
		Message: msg,
	}
	config.DB.Create(&errorLog)
}
