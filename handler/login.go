package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

	//从请求体里获得用户名和密码
	var req struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		Challenge string `json:"challenge"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, 400, "无法解析请求体", err)
		return
	}

	//从数据库里检查这个用户是否存在
	var tmp co.User
	if err := co.DB.First(&tmp, "user_name = ?", req.Username).Error; err != nil {
		util.DbQueryError(c, err, "该用户不存在")
		return
	}

	//检查密码是否正确
	if !AuthChallenge(req.Challenge, req.Password, tmp.Password) {
		util.Error(c, 400, "密码不正确", nil)
		return
	}

	//生成JWT
	token, err := GetJwt(tmp.ID)
	if err != nil {
		util.Error(c, 500, "JWT生成失败", err)
		return
	}

	//将JWT发送给用户
	util.Info(c, 200, "登录成功", gin.H{
		"token":     token,
		"username":  tmp.Username,
		"uid":       tmp.ID,
		"admin":     tmp.Admin,
		"gamename":  tmp.Gamename,
		"email":     tmp.Email,
		"avatar":    tmp.Avatar,
		"telephone": tmp.Telephone,
	})
}
