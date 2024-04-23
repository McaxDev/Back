package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 修改密码的handler
func ChangePwd(c *gin.Context) {

	// 从中间件里读取用户ID
	user, err := BindJwt(c)
	if err != nil {
		util.Error(c, 400, "读取用户信息失败", err)
		return
	}

	// 读取请求体里的密码
	var req struct {
		OldPwd   string `json:"oldPwd"`
		NewPwd   string `json:"newPwd"`
		MailCode string `json:"mailCode"`
	}
	if err := util.BindReq(c, &req); err != nil {
		util.Error(c, 400, "无法读取你的请求体", err)
		return
	}

	// 拒绝没有提供任何验证方式的请求
	if req.MailCode == "" && req.OldPwd == "" {
		util.Error(c, 400, "请至少提供一种验证方式", nil)
		return
	}

	// 提供了邮箱验证码，验证邮箱验证码
	if req.MailCode != "" && !AuthMail(req.MailCode, user.Email) {
		util.Error(c, 400, "邮箱验证码不正确", nil)
		return
	}

	// 提供了原密码，检查原密码是否正确
	if req.OldPwd != "" && util.Encode(req.OldPwd, false) != user.Password {
		util.Error(c, 400, "原密码不正确", nil)
		return
	}

	// 修改密码并保存
	if err := co.DB.Model(user).Update("Password", req.NewPwd).Error; err != nil {
		util.Error(c, 500, "密码修改失败", err)
		return
	}
	util.Info(c, 200, "你的密码修改成功", nil)
}
