package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 用户获取服务器IP地址的handler
func GetIP(c *gin.Context) {

	//查询服务器的IP地址
	var tmp co.IPs
	if err := co.DB.Order("time DESC").First(&tmp).Error; err != nil {
		util.Error(c, 500, "IP地址查询失败", err)
		return
	}

	// 将服务器IP地址发送给用户
	util.Info(c, 200, "IP查询成功", gin.H{
		"IPv4": tmp.Ipv4,
		"IPv6": tmp.Ipv6,
	})
}
