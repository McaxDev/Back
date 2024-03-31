package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

func setAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		util.Error(c, 500, "服务器错误", err)
		return
	}
	username := ReadJwt(c)["name"].(string)
	filename, err := util.RandStr(16)
	if err != nil {
		util.Error(c, 500, "文件名生成失败", err)
		return
	}
	var tmp co.User
	co.DB.First(&tmp, "username = ?", username)
	if tmp.ID != 0 {
		tmp.Avatar = filename
		co.DB.Save(&tmp)
	}
	if err := c.SaveUploadedFile(file, co.Config.AvatarPath); err != nil {
		util.Error(c, 500, "文件保存失败", err)
		return
	}
}
