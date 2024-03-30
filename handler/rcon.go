package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	"github.com/gorcon/rcon"
)

func Rcon(c *gin.Context) {

	//读取用户POST请求的表单数据
	srv, cmd := c.PostForm("server"), c.PostForm("command")
	challenge, hash := c.PostForm("challenge"), c.PostForm("password")

	//判断用户输入的是管理员命令还是普通命令
	otherCmd := true
	for _, item := range co.Config.AllowCmd {
		if cmd == item {
			otherCmd = false
			break
		}
	}

	//阻止密码错误时执行管理员命令
	authed := AuthChallenge(challenge, hash, co.Config.RconPwd)
	if otherCmd && !authed {
		util.Warn(c, 400, "此命令尚未接入", nil)
		return
	}

	//向RCON服务器发送命令
	var conn *rcon.Conn
	var err error
	switch srv {
	case "sc":
		conn, err = exeCmd(co.Config.Port.Sc)
	case "mod":
		conn, err = exeCmd(co.Config.Port.Mod)
	default:
		conn, err = exeCmd(co.Config.Port.Main)
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
	respMap := util.MyMap("info", response)
	util.Info(c, 200, "命令执行成功", respMap)
	return
}

func exeCmd(port string) (*rcon.Conn, error) {
	return rcon.Dial(co.Config.ServerIP+":"+port, co.Config.RconPwd)
}
