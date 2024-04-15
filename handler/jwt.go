package handler

import (
	"errors"
	"fmt"
	"strings"

	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func GetJwt(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
	})
	tokenString, err := token.SignedString([]byte(co.Config.JwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 验证JWT的handler
func AuthJwt(c *gin.Context) {
	Authorization := c.GetHeader("Authorization")
	if !strings.HasPrefix(Authorization, "Bearer ") {
		util.Error(c, 400, "token不合法！", nil)
		return
	}
	RawJwtToken := Authorization[len("Bearer "):]
	JwtToken, err := jwt.Parse(RawJwtToken, keyFunc)
	if err != nil {
		util.Error(c, 400, "token格式不正确！", nil)
		return
	}
	if claims, ok := JwtToken.Claims.(jwt.MapClaims); ok && JwtToken.Valid {
		c.Set("userID", claims["userID"])
	} else {
		util.Error(c, 401, "token身份信息有误！", nil)
		return
	}
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("错误签名方法 %v", token.Header["alg"])
	}
	return []byte(co.Config.JwtKey), nil
}

func ReadJwt(c *gin.Context) (uint, error) {
	jwt, exist := c.Get("userID")
	if !exist {
		return 0, errors.New("无法找到JWT")
	}
	userid, ok := jwt.(float64)
	if !ok {
		return 0, errors.New("对JWT中的用户ID断言失败")
	}
	return uint(userid), nil
}
