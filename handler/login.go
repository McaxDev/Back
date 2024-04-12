package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

	//从请求体里获得用户名和密码
	username, password := c.PostForm("username"), c.PostForm("password")

	//从数据库里检查这个用户是否存在
	var tmp co.User
	err := co.DB.Where("username = ?", username).First(&tmp).Error
	if err != nil {
		util.DbQueryError(c, err, "该用户不存在")
		return
	}

	//检查密码是否正确
	if tmp.Password != util.Encode(password) {
		util.Error(c, 400, "密码不正确", err)
		return
	}

	//生成JWT
	token, err := GetJwt(tmp.ID)
	if err != nil {
		util.Error(c, 500, "JWT生成失败", err)
		return
	}

	//将JWT发送给用户
	util.Info(c, 200, "JWT生成成功", gin.H{"token": token})
}
