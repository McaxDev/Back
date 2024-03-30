package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

func Signup(c *gin.Context) {
	username, password := c.PostForm("username"), c.PostForm("password")
	if err := passwordvalidator.Validate(password, 60.0); err != nil {
		util.Warn(c, 400, "注册失败，密码复杂度不够", err)
		return
	}
	result := co.DB.Where("user_name = ?", username).First(&co.User{})
	if err := result.Error; err == nil {
		util.Warn(c, 403, "该用户已存在", err)
		return
	}
	user := co.User{Username: username, Password: password}
	if err := co.DB.Create(&user).Error; err != nil {
		util.Error(c, 500, "无法创建用户", err)
		return
	}
	util.Info(c, 200, "用户创建成功", nil)
}
