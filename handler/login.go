package handler

import (
	"errors"

	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Login(c *gin.Context) {
	challenge := c.PostForm("challenge")
	hash := c.PostForm("hash")
	username := c.PostForm("username")
	var tmp co.User
	err := co.DB.Where("username = ?", username).First(&tmp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Warn(c, 401, "该用户不存在", err)
		} else {
			util.Error(c, 500, "服务器错误", err)
		}
		return
	}
	if !AuthChallenge(challenge, hash, tmp.Password) {
		util.Warn(c, 403, "密码不正确", nil)
		return
	}
	token, err := GetJwt(int(tmp.ID), tmp.Username, tmp.Admin)
	if err != nil {
		util.Error(c, 500, "JWT生成失败", err)
		return
	}
	tokenMap := map[string]interface{}{"token": token}
	util.Info(c, 200, "JWT生成成功", tokenMap)
}
