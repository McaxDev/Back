package handler

import (
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 自动登录的函数
func AutoLogin(c *gin.Context) {

	// 从中间件jwt读取用户ID
	userID, err := ReadJwt(c)
	if err != nil {
		util.Error(c, 500, "读取JWT失败", err)
		return
	}

	// 读取用户信息
	userdata, err := util.GetUserInfo(userID)
	if err != nil {
		util.Error(c, 500, "用户信息读取失败", err)
		return
	}

	// 将用户信息返回给用户
	util.Info(c, 200, "请求成功", userdata)
}
