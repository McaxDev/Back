package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/McaxDev/Back/config"
)

func Gpt(c *gin.Context) {
	// 获取查询字符串参数
	text := c.Query("text")
	if text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺乏查询字符串参数"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "不合法的temperature值"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "请求失败"})
		return
	}
	gptUrl := "https://api.openai.com/v1/chat/completions"
	gptRequest := strings.NewReader(string(requestBody))
	req, err := http.NewRequest("POST", gptUrl, gptRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "post make fail"})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Config.GptToken)

	// 向GPT发送请求之后向用户发送请求
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", body)
}
