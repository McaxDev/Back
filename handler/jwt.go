package handler

import (
	"fmt"
	"strings"

	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func GetJwt(id int, name string, admin int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"name":  name,
		"admin": admin,
	})
	tokenString, err := token.SignedString(co.Config.JwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Jwt(c *gin.Context) {
	Authorization := c.GetHeader("Authorization")
	if !strings.HasPrefix(Authorization, "Bearer ") {
		util.Warn(c, 400, "token不合法！", nil)
		return
	}
	RawJwtToken := Authorization[len("Bearer "):]
	JwtToken, err := jwt.Parse(RawJwtToken, keyFunc)
	if err != nil {
		util.Warn(c, 400, "token格式不正确！", nil)
		return
	}
	if claims, ok := JwtToken.Claims.(jwt.MapClaims); ok && JwtToken.Valid {
		c.Set("userInfo", claims)
		c.Next()
	} else {
		util.Warn(c, 401, "token身份信息有误！", nil)
		return
	}
}

func ReadJwt(c *gin.Context) map[string]interface{} {
	userInfo, exists := c.Get("userInfo")
	emptyMap := util.MyMap("id", 0, "name", "", "admin", 0)
	if !exists {
		util.Warn(c, 401, "解读JWT失败", nil)
		return emptyMap
	}
	claims, ok := userInfo.(jwt.MapClaims)
	if !ok {
		util.Error(c, 500, "JWT类型断言失败", nil)
		return emptyMap
	}
	return claims
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("错误签名方法 %v", token.Header["alg"])
	}
	return []byte(co.Config.JwtKey), nil
}
