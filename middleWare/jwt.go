package middleWare

import (
	"fmt"
	"strings"

	"github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func GetJwt(id int, name string, admin bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"name":  name,
		"admin": admin,
	})
	tokenString, err := token.SignedString(config.Config.JwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Jwt(c *gin.Context) {
	Authorization := c.GetHeader("Authorization")
	if !strings.HasPrefix(Authorization, "Bearer ") {
		c.AbortWithStatusJSON(400, util.Json("token不合法！", nil))
		return
	}
	RawJwtToken := Authorization[len("Bearer "):]
	JwtToken, err := jwt.Parse(RawJwtToken, keyFunc)
	if err != nil {
		c.AbortWithStatusJSON(401, util.Json("token格式不正确！", nil))
		return
	}
	if claims, ok := JwtToken.Claims.(jwt.MapClaims); ok && JwtToken.Valid {
		c.Set("userInfo", claims)
		c.Next()
	} else {
		c.AbortWithStatusJSON(401, util.Json("token身份信息有误！", nil))
		return
	}
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("错误签名方法 %v", token.Header["alg"])
	}
	return config.Config.JwtKey, nil
}
