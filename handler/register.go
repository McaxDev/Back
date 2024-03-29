package handler

import (
	"github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/entity"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

func Register(c *gin.Context) {
	var newUser entity.User
	if err := c.BindJSON(&newUser); err != nil {
		util.Error(c, 400, "将表单数据绑定到结构体失败", err)
		return
	}
	if err := passwordvalidator.Validate(newUser.UserPas, 60.0); err != nil {
		util.Error(c, 400, "注册失败，密码复杂度不够", err)
		return
	}
	var tmp entity.User
	result := config.DB.Where("user_name = ?", newUser.UserName).First(&tmp)
	if result.Error == nil {
		util.Error(c, 409, "该用户已存在", nil)
		return
	}
	if config.DB.Create(&newUser) != nil {
		util.Error(c, 500, "无法创建用户", nil)
		return
	}
	util.Error(c, 200, "用户创建成功", nil)
}
