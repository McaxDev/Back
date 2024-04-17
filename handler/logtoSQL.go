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

// 将请求的信息存储到数据库的中间件
func LogToSQL(c *gin.Context) {

	// 读取请求体并复制
	reqBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		co.SysLog("ERROR", err.Error())
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	reqBodyStr := string(reqBody)

	//获取请求的持续时间
	startTime := time.Now()
	c.Next()
	duration := time.Since(startTime)

	//从handler里获取错误信息
	midErr, _ := util.ReadMid[error](c, "error")
	errString := util.ErrToStr(midErr)

	//从handler里获取用户ID
	userID, _ := ReadJwt(c)

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
