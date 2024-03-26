package middleWare

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var chals = make(map[string]time.Time)
var chalIte = 0

func GetChallenge(c *gin.Context) {
	chalIteStr := strconv.Itoa(chalIte)
	chals[chalIteStr] = time.Now().Add(time.Minute)
	c.JSON(200, gin.H{"challenge": chalIteStr, "date": chals[chalIteStr]})
	chalIte++
}

func AuthChallenge(realPwd string) gin.HandlerFunc {
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
