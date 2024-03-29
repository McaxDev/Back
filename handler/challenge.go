package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

var chals = make(map[string]time.Time)

func GetChallenge(c *gin.Context) {
	str, err := util.RandStr(16)
	if err != nil {
		util.Error(c, 500, "随机数生成失败", err)
		return
	}
	if _, exists := chals[str]; exists {
		util.Warn(c, 500, "你太幸运了，请重试", nil)
		return
	}
	chals[str] = time.Now().Add(time.Minute)
	data := gin.H{"challenge": str, "date": chals[str]}
	util.Info(c, 200, "获取挑战值成功", data)
}

func AuthChallenge(challenge, hash, password string) bool {
	expiry, exist := chals[challenge]
	if !exist || time.Now().After(expiry) {
		return false
	}
	delete(chals, challenge)
	hashBytes := sha256.Sum256([]byte(challenge + password))
	return hash == hex.EncodeToString(hashBytes[:])
}

func ClearExpiredChallenge() {
	now := time.Now()
	for challenge, expiry := range chals {
		if now.After(expiry) {
			delete(chals, challenge)
		}
	}
}
