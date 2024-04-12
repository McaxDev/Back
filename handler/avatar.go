package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

/*
func SetAvatar(c *gin.Context) {

	//接收用户上传的头像图片
	file, err := c.FormFile("avatar")
	if err != nil {
		util.Error(c, 500, "服务器错误", err)
		return
	}

	//从JWT里获取用户ID
	userID, err := ReadJwt(c)
	if err != nil {
		util.Error(c, 500, "读取用户jwt失败", err)
		return
	}

	//生成随机字符串作为文件名
	filename, err := util.RandStr(16)
	if err != nil {
		util.Error(c, 500, "文件名生成失败", err)
		return
	}

	//从数据库里查找用户
	var tmp co.User
	if err = co.DB.First(&tmp, "user_id = ?", userID).Error; err != nil {
		util.DbQueryError(c, err, "无法找到这个用户")
		return
	}
	if tmp.ID != 0 {
		tmp.Avatar = filename
		co.DB.Save(&tmp)
	}
	if err := c.SaveUploadedFile(file, co.Config.AvatarPath); err != nil {
		util.Error(c, 500, "文件保存失败", err)
		return
	}
}
*/

func GetAvatar(c *gin.Context) {

	//从http请求里获得用户名
	username := c.Query("username")

	//通过用户名在数据库里查找用户
	var tmp co.User
	if err := co.DB.First(&tmp, "user_name = ?", username).Error; err != nil {
		util.DbQueryError(c, err, "无法找到这个用户")
		return
	}

	//将用户头像字符串作为响应返回
	avatarMap := util.Data("avatar", tmp.Avatar)
	util.Info(c, 200, "获取头像成功", avatarMap)
}
