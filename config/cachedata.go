package config

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

type UserData struct {
	score1 int `json:"score1"`
}

func CacheData(server string) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "Cb8b318b",
		DB:       0,
	})
	ctx := context.Background()
	datapath := SrvConf["path"][server] + "world/stats/"
	files, err := os.ReadDir(datapath)
	if err != nil {
		LogFunc("ERROR", err)
		return
	}
	for _, file := range files {
		data, err := os.ReadFile(datapath + file.Name())
		if err != nil {
			LogFunc("ERROR", err)
			continue
		}
	}
}
