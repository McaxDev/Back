package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 修改游戏名的handler
func ChangeGamename(c *gin.Context) {

	// 从jwt里获取用户信息并绑定到结构体对象
	user, err := BindJwt(c)
	if err != nil {
		util.Error(c, 400, "无法读取你的个人信息", err)
		return
	}

	// 保存游戏名到数据库
	if err := co.DB.Model(user).Updates(gin.H{
		"Gamename": c.Query("gamename"),
		"GameAuth": false,
	}).Error; err != nil {
		util.Error(c, 500, "修改游戏名失败", err)
		return
	}

	// 返回应答
	util.Info(c, 200, "修改游戏名成功，请前往验证", nil)
}
