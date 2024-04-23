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

	// 创建代表用户设备的结构体指针
	device := &co.CoinType{Azure: make(chan int, 1), Pearl: make(chan int, 1)}

	// 查询用户余额
	var coin co.AxolotlCoin
	if err := co.DB.First(&coin, "user_id = ?", userID).Error; err != nil {
		util.DbQueryError(c, err, "找不到你的余额记录")
		return
	}

	// 将用户余额填充到结构体里
	device.Azure <- coin.Azure
	device.Pearl <- coin.Pearl

	// 将代表用户设备的结构体指针加入到用户设备列表
	devices := co.CoinPubSub.User[userID]
	co.CoinPubSub.Mu.Lock()
	devices = append(devices, device)
	co.CoinPubSub.Mu.Unlock()

	// 使用循环从通道里读取余额并发送
	for {
		select {
		case azureBalance := <-device.Azure:
			pearlBalance := <-device.Pearl
			balance := map[string]int{"azure": azureBalance, "pearl": pearlBalance}
			c.SSEvent("balance", util.Resp("余额已更新", balance))
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}
