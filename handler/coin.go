package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 使用流传输来实时发送用户的余额
func Coin(c *gin.Context) {

	// 从JWT里读取用户ID
	userID, err := ReadJwt(c)
	if err != nil {
		util.Error(c, 400, "无法读取你的JWT", err)
		return
	}

	// 设置响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// 订阅创建传递余额的通道
	ch, err := co.CoinPubSub.Sub(userID)
	if err != nil {
		util.Error(c, 500, "无法在数据库里查找你的余额", err)
		return
	}
	defer co.CoinPubSub.UnSub(userID)

	// 使用循环从通道里读取余额并发送
	for {
		select {
		case balance := <-ch:
			c.SSEvent("balance", balance)
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}
