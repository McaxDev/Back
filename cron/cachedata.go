package cron

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	co "github.com/McaxDev/Back/config"
	"github.com/go-redis/redis/v8"
)

// 存储用于反序列化JSON数据的结构体
type PlayerData struct {
	Mined   map[string]int `json:"minecraft:mined"`
	Crafted map[string]int `json:"minecraft:crafted"`
	Custom  struct {
		PlayTime       int `json:"minecraft:play_time"`
		DeathCount     int `json:"minecraft:deaths"`
		RideDistance   int `json:"minecraft:horse_one_cm"`
		SwimDistance   int `json:"minecraft:swim_one_cm"`
		WalkDistance   int `json:"minecraft:walk_one_cm"`
		FlyDistance    int `json:"minecraft:fly_one_cm"`
		BoatDistance   int `json:"minecraft:boat_one_cm"`
		SprintDistance int `json:"minecraft:sprint_one_cm"`
		MobKilled      int `json:"minecraft:mob_kills"`
		DamageCaused   int `json:"minecraft:damage_dealt"`
	} `json:"minecraft:custom"`
}

// 将玩家的游戏数据缓存到Redis的函数
func CacheData(server string) error {

	// 读取玩家数据目录并存储到变量
	rootpath := filepath.Join(co.Config.ServerPath[server], "world", "stats")
	datadir, err := os.ReadDir(rootpath)
	if err != nil {
		return err
	}

	// 对数据文件夹里文件进行遍历存储到Redis
	for _, file := range datadir {

		// 将每个数据文件内容读取并汇总到结构体对象
		filename := file.Name()
		path := filepath.Join(rootpath, filename)
		result, err := SummarizeData(path)
		if err != nil {
			co.SysLog("ERROR", "读取玩家游戏数据失败："+err.Error())
			continue
		}

		// 根据配置文件的要求，将每个玩家的数据存储到Redis
		playeruuid := strings.TrimSuffix(filename, ".json")
		for key, value := range *result {
			if err := co.RDB.ZAdd(context.Background(), server+key, &redis.Z{
				Score:  float64(value),
				Member: co.PlayerName[playeruuid],
			}).Err(); err != nil {
				co.SysLog("ERROR", "将玩家数据缓存到Redis失败："+err.Error())
				continue
			}
		}
	}

	return nil
}

// 根据文件路径提取数据并返回汇总后的结构体
func SummarizeData(path string) (*map[string]int, error) {

	// 将文件数据读取到字节切片
	databyte, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// 将字节切片反序列化到结构体对象并返回
	var data PlayerData
	if err := json.Unmarshal(databyte, &data); err != nil {
		return nil, err
	}

	// 将反序列化得到的结构体对象进行汇总返回
	return &map[string]int{
		"Mined":     GetSum(data.Mined),
		"Crafted":   GetSum(data.Crafted),
		"PlayTime":  data.Custom.PlayTime,
		"Death":     data.Custom.DeathCount,
		"Explore":   data.Custom.WalkDistance,
		"Damage":    data.Custom.DamageCaused,
		"MobKilled": data.Custom.MobKilled,
	}, nil
}

// 求和函数
func GetSum(origin map[string]int) (sum int) {
	for _, value := range origin {
		sum += value
	}
	return
}
