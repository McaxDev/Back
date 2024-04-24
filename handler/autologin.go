package handler

import (
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 自动登录的函数
func AutoLogin(c *gin.Context) {

	// 从中间件jwt读取用户ID
	user, err := BindJwt(c)
	if err != nil {
		util.Error(c, 500, "读取JWT失败", err)
		return
	}

	// 将用户信息返回给用户
	util.Info(c, 200, "请求成功", util.GetUserInfo(user))
}
