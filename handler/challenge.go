package handler

import (
	"time"

	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

var challenges = make(map[string]time.Time)

func GetChallenge(c *gin.Context) {
	str, err := util.RandStr(16)
	if err != nil {
		util.Error(c, 500, "随机数生成失败", err)
		return
	}
	if _, exists := challenges[str]; exists {
		util.Error(c, 500, "你太幸运了，请重试", nil)
		return
	}
	challenges[str] = time.Now().Add(time.Minute)
	data := gin.H{"challenge": str, "date": challenges[str]}
	util.Info(c, 200, "获取挑战值成功", data)
}

func AuthChallenge(challenge, hash, password string) bool {
	expiry, exist := challenges[challenge]
	if !exist || time.Now().After(expiry) {
		return false
	}
	delete(challenges, challenge)
	return hash == util.Encode(challenge+password)
}

func ClearExpiredChallenge() {
	now := time.Now()
	for challenge, expiry := range challenges {
		if now.After(expiry) {
			delete(challenges, challenge)
		}
	}
}
