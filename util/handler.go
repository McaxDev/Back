package util

import co "github.com/McaxDev/Back/config"

// 根据用户ID查询用户信息的函数
func GetUserInfo(userid uint) (user co.User, err error) {
	err = co.DB.First(&user, "user_id = ?", userid).Error
	user.Password = "success"
	return
}
