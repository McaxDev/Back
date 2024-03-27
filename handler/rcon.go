package handler

import (
	"github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	"github.com/gorcon/rcon"
)

func Rcon(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.JSON(500, util.Json("失败", "解析表单数据失败"))
		return
	}
	srv, cmd := c.PostForm("server"), c.PostForm("command")
	otherCmd := true
	for _, item := range config.Config.AllowCmd {
		if cmd == item {
			otherCmd = false
			break
		}
	}
	if _, authExist := c.Get("auth"); otherCmd && !authExist {
		c.JSON(400, util.Json("失败", "此命令尚未接入"))
		return
	}
	var conn *rcon.Conn
	var err error
	switch srv {
	case "sc":
		conn, err = rcon.Dial("192.168.50.38:25577", config.Config.RconPwd)
	case "mod":
		conn, err = rcon.Dial("192.168.50.38:25574", config.Config.RconPwd)
	default:
		conn, err = rcon.Dial("192.168.50.38:25575", config.Config.RconPwd)
	}
	if err != nil {
		c.JSON(500, util.Json("失败", "连接RCON服务器失败"))
		return
	}
	defer conn.Close()

	response, err := conn.Execute(cmd)
	if err != nil {
		c.JSON(400, util.Json("失败", "命令执行失败"))
		return
	}
	c.JSON(200, util.Json("失败", response))
	return
}
