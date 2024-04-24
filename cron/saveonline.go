package cron

import "github.com/McaxDev/Back/util"

// 将玩家在线信息定期存储到数据库的函数
func SaveOnline() error {

	// 创建一个存储服务器在线人数的哈希表
	onlineMap := map[string]int{"main": 0, "sc": 0, "mod": 0}

	// 查询服务器信息
	for server := range onlineMap {
		resp, err := util.Status(server)
		if err != nil {
			return err
		}
		onlineMap[server] = int(resp.OnlinePlayers)
	}
	return nil
}
