package handler

import (
	"github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

func Register(c *gin.Context) {
	config.DB.AutoMigrate(&config.User{})
	var newUser config.User
	if c.BindJSON(&newUser) != nil {
		c.JSON(400, util.Json("将表单数据绑定到结构体失败", nil))
		return
	}
	if passwordvalidator.Validate(newUser.UserPas, 60.0) != nil {
		c.JSON(400, util.Json("注册失败，密码复杂度不够", nil))
		return
	}
	var tmp config.User
	result := config.DB.Where("user_name = ?", newUser.UserName).First(&tmp)
	if result.Error == nil {
		c.JSON(409, util.Json("该用户已存在", nil))
		return
	}
	if config.DB.Create(&newUser) != nil {
		c.JSON(500, util.Json("无法创建用户", nil))
		return
	}
	c.JSON(200, util.Json("用户创建成功", nil))
}
