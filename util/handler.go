package util

import (
	"encoding/json"
	"errors"
	"fmt"

	co "github.com/McaxDev/Back/config"
	"github.com/gin-gonic/gin"
)

// 根据用户ID查询用户信息的函数
func GetUserInfo(userid uint) (user co.User, err error) {
	err = co.DB.First(&user, "user_id = ?", userid).Error
	user.Password = "success"
	return
}

// 从中间件读取变量的函数
func ReadMid[T any](c *gin.Context, data string) (T, error) {

	// 检查这个变量的存在性
	variable, exist := c.Get(data)
	if !exist {
		return *new(T), errors.New("中间件里不存在这个变量")
	}

	//对这个变量进行类型断言
	value, ok := variable.(T)
	if !ok {
		return *new(T), fmt.Errorf("对中间件的变量进行类型断言失败：%T", *new(T))
	}
	return value, nil
}

// 将请求体读取到结构体
func BindReq(c *gin.Context, objPtr interface{}) error {

	// 读取未经断言的请求体
	bodyAny, exists := c.Get("reqBody")
	if !exists {
		return errors.New("无法从中间件里读取请求体")
	}

	// 对从中间件读取的请求体进行类型断言
	bodyBytes, ok := bodyAny.([]byte)
	if !ok {
		return errors.New("对来自中间件的请求体类型断言失败")
	}

	// 将数据绑定到结构体
	if err := json.Unmarshal(bodyBytes, objPtr); err != nil {
		return err
	}

	// 如果执行成功，就不返回错误
	return nil
}
