package handler

import (
	"fmt"
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
func AuthChallenge(challenge, hash, password string) bool {
	expiry, exist := Challenges[challenge]
	if !exist {
		fmt.Println("挑战值不存在")
		return false
	}
	if time.Now().After(expiry) {
		fmt.Println("挑战值已过期")
		return false
	}
	// 如果挑战值存在且没有过期，打印出用于比较的哈希值
	calculatedHash := util.Encode(password+challenge, false)
	fmt.Println("真正的哈希值：" + calculatedHash)

	// 删除挑战值
	delete(Challenges, challenge)

	return hash == calculatedHash
}
