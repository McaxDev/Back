package handler

import (
	"time"

	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

var Challenges = make(map[string]time.Time)

func GetChallenge(c *gin.Context) {
	str := util.RandStr(16)
	if _, exists := Challenges[str]; exists {
		util.Error(c, 500, "你太幸运了，请重试", nil)
		return
	}
	Challenges[str] = time.Now().Add(time.Minute)
	data := gin.H{"challenge": str, "date": Challenges[str]}
	util.Info(c, 200, "获取挑战值成功", data)
}

/*
	func AuthChallenge(challenge, hash, password string) bool {
		expiry, exist := Challenges[challenge]
		if !exist || time.Now().After(expiry) {
			return false
		}
		delete(Challenges, challenge)
		fmt.Println("真正的哈希值" + util.Encode(password+challenge, false))
		return hash == util.Encode(password+challenge, false)
	}
*/

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
