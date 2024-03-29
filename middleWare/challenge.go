package middleWare

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

var chals = make(map[string]time.Time)
var chalIte = 0

func GetChallenge(c *gin.Context) {
	str, err := util.RandStr(16)
	if err != nil {
		util.Error(c, 500, "随机数生成失败", nil, err)
		return
	}
	if _, exists := chals[str]; exists {
		util.Error(c, 500, "your luck good", nil, err)
		return
	}
	chals[str] = time.Now().Add(time.Minute)
	data := gin.H{"challenge": str, "date": chals[str]}

}

func Challengex(realPwd string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userChal, userHash := c.PostForm("challenge"), c.PostForm("password")
		if time.Now().Before(chals[userChal]) {
			delete(chals, userChal)
			hashBytes := sha256.Sum256([]byte(realPwd + userChal))
			if userHash == hex.EncodeToString(hashBytes[:]) {
				c.Set("auth", true)
			}
		}
		c.Next()
	}
}

func ClearExpiredChallenge() {
	now := time.Now()
	for challenge, expiry := range chals {
		if now.After(expiry) {
			delete(chals, challenge)
		}
	}
}
