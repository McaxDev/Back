package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/McaxDev/Back/config"
	"github.com/gin-gonic/gin"
	"github.com/gorcon/rcon"
)

func Listenrcon(c *gin.Context) {
	switch c.Request.Method {

	// 对GET请求返回挑战随机数
	case http.MethodGet:
		randBytes := make([]byte, 16)
		if _, err := rand.Read(randBytes); err != nil {
			c.String(http.StatusInternalServerError, "随机数生成失败")
			return
		}
		randStr := hex.EncodeToString(randBytes)
		challenges[chalIte] = randStr
		chalIte = (chalIte + 1) % chalVolume
		c.String(http.StatusOK, randStr)

	// 对POST请求验证并执行命令
	case http.MethodPost:

		// 解析表单数据
		if err := c.Request.ParseForm(); err != nil {
			c.String(http.StatusBadRequest, "解析表单数据失败")
			return
		}
		srv := c.PostForm("server")
		cmd := c.PostForm("command")
		chal := c.PostForm("challenge")
		hash := c.PostForm("hash")

		// 验证挑战值是否存在
		chalExist := false
		for _, tempChal := range challenges {
			if tempChal == chal {
				chalExist = true
				break
			}
		}
		if !chalExist {
			c.String(http.StatusBadRequest, "挑战值不存在")
			return
		}

		// 验证哈希值是否正确
		var RCONpwd string
		if temp, ok := config.Conf["RCONpwd"].(string); ok {
			RCONpwd = temp
		}
		hashBytes := sha256.Sum256([]byte(RCONpwd + chal))
		realhash := hex.EncodeToString(hashBytes[:])
		if realhash != hash {
			c.String(http.StatusBadRequest, "密码错误")
			return
		}

		// 执行RCON命令
		var conn *rcon.Conn
		var err error
		switch srv {
		case "sc":
			conn, err = rcon.Dial("192.168.50.38:25577", RCONpwd)
		case "mod":
			conn, err = rcon.Dial("192.168.50.38:25574", RCONpwd)
		default:
			conn, err = rcon.Dial("192.168.50.38:25575", RCONpwd)
		}
		if err != nil {
			c.String(http.StatusInternalServerError, "连接RCON服务器失败")
			return
		}
		defer conn.Close()

		response, err := conn.Execute(cmd)
		if err != nil {
			c.String(http.StatusInternalServerError, "命令执行失败")
			return
		}
		c.String(http.StatusOK, response)
	}
}
