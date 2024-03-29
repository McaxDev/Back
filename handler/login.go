package handler

import (
	"errors"

	"github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/entity"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Login(c *gin.Context) {
	var user, tmp entity.User
	if err := c.ShouldBind(&user); err != nil {
		util.Error(c, 400, "请求表单数据有误", err)
		return
	}
	result := config.DB.Where("username = ?", user.UserName).First(&tmp)
	if err := result.Error; err != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			util.Error(c, 401, "该用户不存在", err)
		} else {
			util.Error(c, 500, "服务器错误", err)
		}
		return
	}
	if user.UserPas != tmp.UserPas {
		util.Error(c, 401, "密码不正确", nil)
		return
	}
}
