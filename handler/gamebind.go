package handler

import (
	"fmt"
	"time"

	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

type bindStruct struct {
	Gamename string
	Expire   time.Time
}

var bindcodes = make(map[string]bindStruct)

func GameBindCode(c *gin.Context) {
	srv, gamename := c.PostForm("server"), c.PostForm("gamename")
	bindcode := util.RandStr(6)
	bindcodes[bindcode] = bindStruct{
		Gamename: gamename,
		Expire:   time.Now().Add(10 * time.Minute),
	}
	command := fmt.Sprintf("tell %s 你的验证码是 %s", gamename, bindcode)
	_, err := util.Rcon(srv, command)
	if err != nil {
		util.Error(c, 500, "验证码无法送达MC服务器", err)
		return
	}
	util.Info(c, 200, "验证码已送达MC服务器，十分钟内有效", nil)
}

func AuthBindCode(c *gin.Context) {
	authcode := c.PostForm("authcode")
	authObj, exist := bindcodes[authcode]
	if !exist || time.Now().After(authObj.Expire) {
		util.Error(c, 400, "验证码无效或已过期", nil)
		return
	}
	username, _ := ReadJwt(c)
	err := co.DB.Model(&co.User{}).Where("username = ?", username).Update("gamename", authObj.Gamename).Error
	if err != nil {
		util.Error(c, 500, "数据库查询用户失败", err)
		return
	}
	util.Info(c, 200, "你账号绑定的玩家已更新", nil)
}
