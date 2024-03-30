package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

func GetIP(c *gin.Context) {
	var tmp co.IPs
	if err := co.DB.Order("time DESC").First(&tmp).Error; err != nil {
		util.Error(c, 500, "IP地址查询失败", err)
		return
	}
	IpMap := util.MyMap("IPv4", tmp.Ipv4, "IPv6", tmp.Ipv6)
	util.Info(c, 200, "IP查询成功", IpMap)
}
