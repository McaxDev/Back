package handler

import (
	"time"

	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

var Challenges = make(map[string]time.Time)

func GetChallenge(c *gin.Context) {
	str, err := util.RandStr(16)
	if err != nil {
		util.Error(c, 500, "随机数生成失败", err)
		return
	}
	if _, exists := Challenges[str]; exists {
		util.Error(c, 500, "你太幸运了，请重试", nil)
		return
	}
	Challenges[str] = time.Now().Add(time.Minute)
	data := gin.H{"challenge": str, "date": Challenges[str]}
	util.Info(c, 200, "获取挑战值成功", data)
}

func AuthChallenge(challenge, hash, password string) bool {
	expiry, exist := Challenges[challenge]
	if !exist || time.Now().After(expiry) {
		return false
	}
	delete(Challenges, challenge)
	return hash == util.Encode(challenge+password, true)
}
