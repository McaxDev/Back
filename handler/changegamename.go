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

	// 从请求头获取游戏名
	gamename := c.Query("gamename")

	// 保存游戏名到数据库
	result := co.DB.Model(user).Updates(gin.H{"Gamename": gamename, "GameAuth": false})
	if err := result.Error; err != nil {
		util.Error(c, 500, "修改游戏名不正确", err)
		return
	}

	// 返回应答
	util.Info(c, 200, "修改游戏名正确，请前往验证", nil)
}
