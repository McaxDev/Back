package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	"github.com/mcstatus-io/mcutil/v3"
)

// 查询服务器状态的handler
func Status(c *gin.Context) {

	// 创建一个超时的context
	ctx, canc := util.Timeout(5)
	defer canc()

	// 将用户的请求体绑定结构体
	var req struct {
		Server string `json:"server"`
	}
	if err := util.BindReq(c, &req); err != nil {
		util.Error(c, 500, "无法读取你的请求体", err)
		return
	}

	// 执行查询操作
	port := uint16(co.Config.ServerPort[req.Server])
	resp, err := mcutil.FullQuery(ctx, co.Config.ServerIP, port)
	if err != nil {
		util.Error(c, 500, "查询失败", err)
		return
	}

	// 将查询完成的响应发送给用户
	util.Info(c, 200, "查询成功", resp)
}
