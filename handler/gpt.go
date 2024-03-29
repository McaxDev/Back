package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

func Gpt(c *gin.Context) {
	// 获取查询字符串参数
	text := c.Query("text")
	if text == "" {
		util.Error(c, 400, "缺乏查询字符串参数", err)
		return
	}
	model := "gpt-3.5-turbo"
	if c.Query("model") == "4" {
		model = "gpt-4"
	}
	temperature := 0.7
	if inputed := c.Query("temperature"); inputed != "" {
		temp, err := strconv.ParseFloat(inputed, 64)
		if err != nil {
			util.Error(c, 400, "不合法的temperature值", err)
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
		util.Error(c, 500, "请求失败", err)
		return
	}
	gptUrl := "https://api.openai.com/v1/chat/completions"
	gptRequest := strings.NewReader(string(requestBody))
	req, err := http.NewRequest("POST", gptUrl, gptRequest)
	if err != nil {
		util.Error(c, 500, "请求发送失败", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer ")

	// 向GPT发送请求之后向用户发送请求
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		util.Error(c, 500, err.Error(), err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.Error(c, 500, err.Error(), err)
		return
	}
	c.Data(http.StatusOK, "application/json", body)
}
