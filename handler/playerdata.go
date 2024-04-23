package handler

import (
	"path/filepath"

	co "github.com/McaxDev/Back/config"
	cr "github.com/McaxDev/Back/cron"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 获取玩家游戏数据的handler
func PlayerData(c *gin.Context) {

	// 从中间件JWT里读取玩家信息
	user, err := BindJwt(c)
	if err != nil {
		util.Error(c, 400, "无法读取你的个人信息", err)
		return
	}

	// 根据玩家名和服务器信息得到文件路径
	server := c.Query("server")
	filename := co.PlayerUUID[user.Gamename] + ".json"
	path := filepath.Join(co.Config.ServerPath[server], filename)

	// 根据文件名提取文件中的数据
	data, err := cr.SummarizeData(path)
	if err != nil {
		util.Error(c, 500, "无法找到你的玩家数据", err)
		return
	}

	// 将玩家数据返回给客户端
	util.Info(c, 200, "查询成功", *data)
}
