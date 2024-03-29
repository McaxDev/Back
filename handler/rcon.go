package handler

import (
	"github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	"github.com/gorcon/rcon"
)

func Rcon(c *gin.Context) {

	//读取用户POST请求的表单数据
	if err := c.Request.ParseForm(); err != nil {
		util.Error(c, 500, "解析表单数据失败", err)
		return
	}
	srv, cmd := c.PostForm("server"), c.PostForm("command")
	challenge, hash := c.PostForm("challenge"), c.PostForm("password")

	//判断用户输入的是管理员命令还是普通命令
	otherCmd := true
	for _, item := range config.Config.AllowCmd {
		if cmd == item {
			otherCmd = false
			break
		}
	}

	//阻止密码错误时执行管理员命令
	authed := AuthChallenge(challenge, hash, config.Config.RconPwd)
	if otherCmd && !authed {
		util.Warn(c, 400, "此命令尚未接入", nil)
		return
	}

	//向RCON服务器发送命令
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
		util.Error(c, 500, "连接RCON服务器失败", err)
		return
	}
	defer conn.Close()

	response, err := conn.Execute(cmd)
	if err != nil {
		util.Error(c, 400, "命令执行失败", err)
		return
	}
	respMap := map[string]interface{}{"info": response}
	util.Info(c, 200, "命令执行成功", respMap)
	return
}
