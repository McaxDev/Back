package handler

import (
	"fmt"
	"strings"

	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// 通过用户id生成jwt
func GetJwt(userID uint) (string, error) {

	// 生成jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
	})

	// 对jwt使用密钥进行签名
	tokenString, err := token.SignedString([]byte(co.Config.JwtKey))
	if err != nil {
		return "", err
	}

	// 返回jwt
	return tokenString, nil
}

// 验证JWT的handler
func AuthJwt(c *gin.Context) {

	// 从请求头获取jwt
	Authorization := c.GetHeader("Authorization")

	// 检查jwt是否以Bearer开头
	if !strings.HasPrefix(Authorization, "Bearer ") {
		util.Error(c, 400, "token不合法！", nil)
		return
	}

	// 解析jwt的格式
	RawJwtToken := Authorization[len("Bearer "):]
	JwtToken, err := jwt.Parse(RawJwtToken, keyFunc)
	if err != nil {
		util.Error(c, 400, "token格式不正确！", nil)
		return
	}

	// 检查jwt是否通过
	claims, ok := JwtToken.Claims.(jwt.MapClaims)
	if !ok || !JwtToken.Valid {
		util.Error(c, 401, "token身份信息有误！", nil)
		return
	}

	// 将jwt传递给后续的业务逻辑函数
	c.Set("userID", claims["userID"])
}

// 通过请求里的jwt读取用户ID的函数
func ReadJwt(c *gin.Context) (uint, error) {

	// 从中间件里读取用户id
	userID, err := util.ReadMid[float64](c, "userID")
	if err != nil {
		return 0, err
	}

	// 将userid转换为uint并返回给用户
	return uint(userID), nil
}

// 通过请求里的jwt将用户数据绑定到结构体对象
func BindJwt(c *gin.Context, preloads ...string) (*co.User, error) {

	// 从请求里读取用户ID
	userID, err := ReadJwt(c)
	if err != nil {
		return nil, err
	}

	// 根据用户ID在数据库里查找用户并返回
	var user co.User
	var result = co.DB.Where("user_id = ?", userID)

	//预加载从表的数据并返回
	for _, preload := range preloads {
		result = result.Preload(preload)
	}
	result = result.First(&user)

	// 返回数据和可能的错误
	if err := result.Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// 对生成jwt进行签名的函数
func keyFunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("错误签名方法 %v", token.Header["alg"])
	}
	return []byte(co.Config.JwtKey), nil
}
