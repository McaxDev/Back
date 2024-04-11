package handler

import (
	"bytes"
	"io"
	"log"
	"time"

	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

func LogToSQL(c *gin.Context) {
	startTime := time.Now()
	c.Next()
	duration := time.Since(startTime)
	errString := ""
	if UnknownErr, exist := c.Get("error"); exist {
		if err, ok := UnknownErr.(error); ok {
			errString = util.ErrToStr(err)
		}
	}
	username, _ := ReadJwt(c)
	reqBody, err := io.ReadAll(c.Request.Body)
	reqBodyStr := string(reqBody)
	if err != nil {
		co.ConsoleLog("ERROR", err)
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	DBlog := co.Log{
		Username: username,
		Time:     time.Now(),
		Status:   c.Writer.Status(),
		Error:    errString,
		Duration: duration,
		Method:   c.Request.Method,
		Path:     c.Request.URL.Path,
		Source:   c.ClientIP(),
		ReqBody:  reqBodyStr,
	}
	if dbErr := co.DB.Create(&DBlog).Error; dbErr != nil {
		log.Println("将日志存储到数据库失败：" + dbErr.Error())
	}
}
