package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

func Signup(c *gin.Context) {

	//从表单数据获取用户名，密码，邮箱，验证码， 游戏名
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Authcode string `json:"mailcode"`
		Gamename string `json:"gamename"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, 400, "无法解析请求体", err)
		return
	}

	//验证邮箱验证码是否通过
	if !AuthMail(req.Authcode, req.Email) {
		util.Error(c, 400, "邮箱验证码不正确！", nil)
		return
	}

	//检查密码复杂度是否足够
	if err := passwordvalidator.Validate(req.Password, 60.0); err != nil {
		util.Error(c, 400, "注册失败，密码复杂度不够", err)
		return
	}

	//检查此用户是否已经存在
	var user co.User
	if err := co.DB.First(&user, "user_name = ?", req.Username).Error; err == nil {
		util.Error(c, 403, "该用户已存在", err)
		return
	}

	//将用户信息存储到数据库
	user.Username, user.Password = req.Username, req.Password
	user.Email, user.Gamename = req.Email, req.Gamename
	if err := co.DB.Create(&user).Error; err != nil {
		util.Error(c, 500, "无法创建用户", err)
		return
	}

	//生成JWT
	token, err := GetJwt(user.ID)
	if err != nil {
		util.Error(c, 500, "用户创建成功，但JWT生成失败", err)
		return
	}

	//将JWT发送给用户
	util.Info(c, 200, "用户创建成功", gin.H{"token": token})
}
