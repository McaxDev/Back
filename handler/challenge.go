package handler

import (
	"time"

	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 存储挑战值的数据结构
var Challenges = make(map[string]time.Time)

// 用户获取挑战值
func GetChallenge(c *gin.Context) {

	// 生成并存储挑战值
	str := util.RandStr(16)
	Challenges[str] = time.Now().Add(time.Minute)

	// 将挑战值发送给用户
	util.Info(c, 200, "获取挑战值成功", gin.H{
		"challenge": str,
		"date":      Challenges[str],
	})
}

// 验证挑战值加密的密码是否正确
func AuthChallenge(challenge, hash, password string) bool {

	// 检查挑战值是否过期
	if time.Now().After(Challenges[challenge]) {
		return false
	}
	delete(Challenges, challenge)

	// 检查密码是否正确
	return hash == util.Encode(password+challenge, false)
}
