package cron

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/mcstatus-io/mcutil/v3"
)

// 将玩家在线信息定期存储到数据库的函数
func SaveOnline() error {

	// 创建一个存储服务器在线人数的哈希表
	onlineMap := map[string]int{"main": 0, "sc": 0, "mod": 0}

	// 查询服务器信息
	for server := range onlineMap {
		resp, err := util.Status(server, mcutil.BasicQuery)
		if err != nil {
			return err
		}
		onlineMap[server] = int(resp.OnlinePlayers)
	}

	// 将玩家信息插入数据库
	if err := co.DB.Create(&co.Online{
		Main: onlineMap["main"],
		Sc:   onlineMap["sc"],
		Mod:  onlineMap["mod"],
	}).Error; err != nil {
		return err
	}
	return nil
}
