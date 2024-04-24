package handler

import (
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 查询服务器状态的handler
func Status(c *gin.Context) {

	// 将用户的请求体绑定结构体
	server := c.Query("server")

	// 执行查询操作
	resp, err := util.Status(server)
	if err != nil {
		util.Error(c, 500, "查询服务器信息失败", err)
		return
	}

	// 将查询完成的响应发送给用户
	util.Info(c, 200, "查询成功", resp)
}
