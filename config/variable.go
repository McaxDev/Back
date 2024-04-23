package config

import "github.com/go-redis/redis/v8"

// 存储玩家和对应UUID的哈希表
var PlayerUUID = make(map[string]string)
var PlayerName = make(map[string]string)

// 操作Redis的全局变量
var RDB *redis.Client
