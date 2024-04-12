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

	//获取请求的持续时间
	startTime := time.Now()
	c.Next()
	duration := time.Since(startTime)

	//从handler里获取错误信息
	errString := ""
	if UnknownErr, exist := c.Get("error"); exist {
		if err, ok := UnknownErr.(error); ok {
			errString = util.ErrToStr(err)
		}
	}

	//从handler里获取用户ID
	userID, _ := ReadJwt(c)

	//获取请求体字符串
	reqBody, err := io.ReadAll(c.Request.Body)
	reqBodyStr := string(reqBody)
	if err != nil {
		co.ConsoleLog("ERROR", err)
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))

	//将日志条目存储到数据库
	DBlog := co.Log{
		UserID:   userID,
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
