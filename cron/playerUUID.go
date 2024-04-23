package cron

import (
	"encoding/json"
	"os"
	"path/filepath"

	co "github.com/McaxDev/Back/config"
)

// 用于反序列化JSON的结构体
type playerUUID struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

// 将玩家名与UUID缓存到哈希表的函数
func CachePlayerUUID(server string) {

	// 将文件数据读取到字节切片
	usercache := filepath.Join(co.Config.ServerPath[server], "usercache.json")
	data, err := os.ReadFile(usercache)
	if err != nil {
		co.SysLog("ERROR", err.Error())
		return
	}

	// 将字节切片里的JSON数据反序列化到结构体切片
	var players []playerUUID
	if err := json.Unmarshal(data, &players); err != nil {
		co.SysLog("ERROR", err.Error())
		return
	}

	// 使用循环将玩家名和玩家对应的UUID存储到双向哈希表
	for _, player := range players {
		co.PlayerUUID[player.Name] = player.UUID
		co.PlayerName[player.UUID] = player.Name
	}
}
