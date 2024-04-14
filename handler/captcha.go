package handler

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/McaxDev/Back/util"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/sjson"
)

// 获取captcha验证码
func GetCaptcha(c *gin.Context) {
	id := captcha.New()
	c.Header("Content-Type", "image/png")
	c.Header("X-Captcha-Id", id)
	err := captcha.WriteImage(c.Writer, id, captcha.StdWidth, captcha.StdHeight)
	if err != nil {
		util.Error(c, 500, "验证码绘制失败", err)
		return
	}
}

// 验证captcha验证码
func Captcha(c *gin.Context) {

	//从请求头里读取验证码值和编号
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		util.Error(c, 400, "无法读取请求体", err)
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	//将用户的请求体绑定到服务器映射
	var jsonData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &jsonData); err != nil {
		util.Error(c, 400, "你的请求体无法绑定到服务器映射", err)
		return
	}

	//对captchaID进行类型断言
	captchaID, IDOK := jsonData["captchaID"].(string)
	captchaValue, ValueOK := jsonData["captchaValue"].(string)
	if !IDOK || !ValueOK {
		util.Error(c, 400, "缺少必要的验证码信息", nil)
		return
	}

	//验证captcha验证码
	if !captcha.VerifyString(captchaID, captchaValue) {
		util.Error(c, 400, "验证码不正确", nil)
		return
	}

	//阉割后续逻辑不需要的内容
	modifiedBody, _ := sjson.Delete(string(bodyBytes), "captchaID")
	modifiedBody, _ = sjson.Delete(modifiedBody, "captchaValue")
	c.Request.Body = io.NopCloser(bytes.NewBufferString(modifiedBody))
}
