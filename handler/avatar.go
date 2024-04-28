package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 接受用户上传头像
func SetAvatar(c *gin.Context) {

	// 根据中间件的JWT获取用户身份
	user, err := BindJwt(c)
	if err != nil {
		util.Error(c, 400, "无法获取你的用户信息", err)
		return
	}

	// 将文件上传到文件系统
	innerPath, err := util.UploadFile(c, "avatar", "avatar")
	if err != nil {
		util.Error(c, 500, "文件上传失败", err)
		return
	}

	// 将文件路径存储到历史头像库
	if err := co.DB.Create(&co.UserAvatar{
		UserID: user.ID,
		Path:   innerPath,
	}).Error; err != nil {
		util.Error(c, 500, "无法将你的头像上传到历史头像库", err)
		return
	}

	// 将头像路径返回给用户
	util.Info(c, 200, "上传成功", innerPath)
}

// 获得特定用户的历史头像
func GetAvatar(c *gin.Context) {

	// 从中间件JWT里获取用户身份
	user, err := BindJwt(c, "UserAvatar")
	if err != nil {
		util.Error(c, 400, "无法读取你的用户信息", err)
		return
	}

	// 将用户的历史头像文件名传输给客户端
	util.Info(c, 200, "查询成功", user.UsedAvatar)
}
