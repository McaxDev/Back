package handler

import (
	"time"

	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 存储绑定验证码的哈希表
var bindcodes = make(map[string]bindStruct)

// 上面哈希表的键值对的值
type bindStruct struct {
	Gamename string
	Expire   time.Time
}

// 接收游戏服务器的绑定请求，并发送验证码
func GameBindCode(c *gin.Context) {

	// 将请求体绑定到结构体对象
	gamename := c.Query("gamename")

	// 生成用于绑定验证的验证码，并存储到哈希表
	bindcode := util.RandStr(6)
	bindcodes[bindcode] = bindStruct{
		Gamename: gamename,
		Expire:   time.Now().Add(time.Minute),
	}

	// 将绑定码返回给Minecraft服务器
	util.Info(c, 200, "验证码发送成功，请于10分钟内使用", bindcode)
}

// 验证绑定的handler
func AuthBindCode(c *gin.Context) {

	// 从jwt里获取用户ID
	userID, err := ReadJwt(c)
	if err != nil {
		util.Error(c, 500, "读取用户信息失败", err)
		return
	}

	// 检查用户是否已经通过验证
	var user co.User
	if err := co.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		util.DbQueryError(c, err, "无法找到这个用户")
		return
	}
	if user.GameAuth {
		util.Error(c, 400, "这个用户已经认证过了", err)
		return
	}

	// 从用户的请求里获取验证码
	authcode := c.Query("authcode")

	// 检查绑定验证码是否存在或对应
	bindStru, exist := bindcodes[authcode]
	if !exist || bindStru.Expire.Before(time.Now()) {
		util.Error(c, 400, "验证码无效或过期", nil)
		return
	}

	// 完成绑定
	if err := co.DB.Model(&user).Update("GameAuth", true).Error; err != nil {
		util.Error(c, 500, "绑定失败，系统错误", err)
		return
	}
	util.Info(c, 200, "绑定成功", nil)
}
