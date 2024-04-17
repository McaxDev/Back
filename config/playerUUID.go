package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var PlayerUUID = make(map[string]string)

type playerUUID struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

func CachePlayerUUID(server string) {
	usercache := filepath.Join(Config.ServerPath[server], "usercache.json")
	data, err := os.ReadFile(usercache)
	if err != nil {
		SysLog("ERROR", err.Error())
		return
	}
	var players []playerUUID
	if err := json.Unmarshal(data, &players); err != nil {
		SysLog("ERROR", err.Error())
		return
	}
	for _, player := range players {
		PlayerUUID[player.Name] = player.UUID
	}
}
