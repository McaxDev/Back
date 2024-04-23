package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 修改用户邮箱的handler
func ChangeMail(c *gin.Context) {

	// 从请求体读取用户信息并绑定到结构体对象
	user, err := BindJwt(c)
	if err != nil {
		util.Error(c, 400, "无法读取你的个人信息", err)
		return
	}

	// 从请求体里获取旧邮箱码，新邮箱和新邮箱码
	var req struct {
		OidCode string `json:"oldAuthCode"`
		NewMail string `json:"newMail"`
		NewCode string `json:"newAuthCode"`
	}
	if err := util.BindReq(c, &req); err != nil {
		util.Error(c, 400, "无法读取你的请求体", err)
		return
	}

	// 验证旧的邮箱验证码是否正确
	if !AuthMail(req.OidCode, user.Email) {
		util.Error(c, 400, "旧的邮箱验证码不正确", nil)
		return
	}

	// 验证新的邮箱验证码是否正确
	if !AuthMail(req.NewCode, req.NewMail) {
		util.Error(c, 400, "新的邮箱验证码不正确", nil)
		return
	}

	// 修改用户的邮箱为新的邮箱
	if err := co.DB.Model(user).Update("Email", req.NewMail).Error; err != nil {
		util.Error(c, 500, "修改邮箱失败", err)
		return
	}
	util.Info(c, 200, "邮箱修改成功", nil)
}
