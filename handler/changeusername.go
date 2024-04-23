package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 修改用户名的handler
func ChangeUsername(c *gin.Context) {

	// 从中间件里读取并绑定用户到结构体对象
	user, err := BindJwt(c)
	if err != nil {
		util.Error(c, 400, "无法读取你的个人信息", err)
		return
	}

	// 从请求头里读取用户名
	newname := c.Query("newname")

	// 检查这个用户名是否已经有人使用
	if err := co.DB.First(&co.User{}, "user_name = ?", newname).Error; err == nil {
		util.Error(c, 400, "你选择的用户名已经有人使用", err)
		return
	}

	// 修改用户的用户名
	if err := co.DB.Model(user).Update("Username", newname).Error; err != nil {
		util.Error(c, 500, "修改用户名失败", err)
		return
	}
	util.Info(c, 200, "修改用户名成功", nil)
}
