package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

func Gpt(c *gin.Context) {
	// 获取查询字符串参数
	text := c.Query("text")
	if text == "" {
		util.Error(c, 400, "缺乏查询字符串参数", nil)
		return
	}
	//根据查询字符串判断使用的模型
	model := "gpt-3.5-turbo"
	if c.Query("model") == "4" {
		model = "gpt-4"
	}
	//根据查询字符串设置GPT temperature
	temperature := 0.7
	if inputed := c.Query("temperature"); inputed != "" {
		temp, err := strconv.ParseFloat(inputed, 64)
		if err != nil || temp >= 1.0 || temp <= 0.0 {
			util.Error(c, 400, "不合法的temperature值", nil)
			return
		}
		temperature = temp
	}

	// 制作向GPT发送的请求
	requestBody, err := json.Marshal(map[string]interface{}{
		"model":       model,
		"messages":    []map[string]string{{"role": "user", "content": text}},
		"temperature": temperature,
	})
	if err != nil {
		util.Error(c, 500, "请求创建失败", err)
		return
	}
	gptUrl := "https://api.openai.com/v1/chat/completions"
	gptRequest := strings.NewReader(string(requestBody))
	req, err := http.NewRequest("POST", gptUrl, gptRequest)
	if err != nil {
		util.Error(c, 500, "请求加载失败", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Config.GptToken)

	// 向GPT发送请求之后向用户发送请求
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		util.Error(c, 500, "向OpenAI发送请求失败", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.Error(c, 500, "读取响应失败", err)
		return
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		util.Error(c, 500, "解析响应失败", err)
		return
	}
	util.Info(c, 200, "请求成功", data)
}
