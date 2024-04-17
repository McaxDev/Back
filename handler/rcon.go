package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

func Rcon(c *gin.Context) {

	//将用户的请求体绑定到结构体
	var req struct {
		Server    string `json:"server"`
		Command   string `json:"command"`
		Challenge string `json:"challenge"`
		Hash      string `json:"hash"`
	}
	if err := util.BindReq(c, &req); err != nil {
		util.Error(c, 400, "无法读取你的请求体", err)
		return
	}

	// 阻止密码错误时执行管理员命令
	allowedCmd := util.Contain(co.SrvInfo.AllowCmd, req.Command)
	authed := AuthChallenge(req.Challenge, req.Hash, co.Config.RconPwd)
	if !allowedCmd && !authed {
		util.Error(c, 400, "此命令尚未接入", nil)
		return
	}

	// 向RCON服务器发送命令
	response, err := util.Rcon(req.Server, req.Command)
	if err != nil {
		util.Error(c, 500, "命令发送失败", err)
		return
	}

	// 将命令的结果返回给客户端
	util.Info(c, 200, "命令执行成功", gin.H{"info": response})
	return
}
