package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/McaxDev/Back/config"
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

type Reqres struct {
	Index int
	Id    string
	Pid   string
	Type  string
	Info
}

type Info struct {
	Time     time.Time
	User     string
	UserID   string
	GptModel string
	Content  string
}

func AskGpt(c *gin.Context) {

	//从请求体读取用户的问题，模型，会话
	var reqbody Reqres
	if err := c.BindJSON(&reqbody); err != nil {
		util.Error(c, 400, "无法解析你发送的请求", err)
		return
	}
	question, model, thread := reqbody.Info.Content, reqbody.GptModel, reqbody.Pid

	//从JWT中间件读取用户的ID
	userID, err := ReadJwt(c)
	if err != nil {
		util.Error(c, 500, "JWT读取失败", err)
		return
	}

	//如果使用已有的会话，检查会话id是否存在
	var session co.GptThread
	if thread != "" {
		if err := co.DB.First(&session, "thread = ? AND user_id = ?", thread, userID).Error; err != nil {
			util.DbQueryError(c, err, "找不到对应的会话")
			return
		}
	}

	//检查用户余额是否充足
	var tmp co.AxolotlCoin
	if err = config.DB.First(&tmp, "user_id = ?", userID).Error; err != nil {
		util.DbQueryError(c, err, "找不到对应的用户")
		return
	}
	if (model == "4" && tmp.Azure < 2) || tmp.Pearl < 2 {
		util.Error(c, 400, "你的余额不足，无法提问", nil)
		return
	}

	//创建向GPT发送的http请求
	reqBody := fmt.Sprintf(`{"role": "user", "content": "%s"}`, question)
	url := fmt.Sprintf("https://api.openai.com/v1/threads/%s/messages", thread)
	if thread == "" {
		url = "https://api.openai.com/v1/threads"
		reqBody = fmt.Sprintf(`{"messages": [%s]}`, reqBody)
	}
	req, err := gptRequest(url, reqBody)
	if err != nil {
		util.Error(c, 500, "向GPT发送的请求创建失败", err)
		return
	}

	//向GPT发送http请求
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		util.Error(c, 500, "向GPT发送请求失败", err)
		return
	}
	defer res.Body.Close()

	//读取响应体
	body, err := io.ReadAll(res.Body)
	if err != nil {
		util.Error(c, 500, "响应体读取失败", err)
		return
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		util.Error(c, 500, "JSON反序列化失败", err)
		return
	}

	//对于创建会话，提取新创建会话的id
	if thread == "" {
		if idValue, ok := data["id"].(string); !ok {
			util.Error(c, 500, "对创建会话thread的响应体类型断言失败", nil)
			return
		} else {
			thread = idValue
		}
	}

	//确定用户选择的GPT模型
	var asstid string
	asstStruct := config.Config.AssistantID
	switch model {
	case "GPT4T":
		asstid = asstStruct.Gpt4
	case "HELPER":
		asstid = asstStruct.Axo
	default:
		asstid = asstStruct.Gpt3
	}

	//创建让GPT处理会话中的问题的请求
	url = fmt.Sprintf("https://api.openai.com/v1/threads/%s/runs", thread)
	reqBody = fmt.Sprintf(`{"assistant_id": "%s"}`, asstid)
	req, err = gptRequest(url, reqBody)
	if err != nil {
		util.Error(c, 500, "执行GPT的http请求创建失败", err)
		return
	}

	//发送让GPT处理会话问题的请求
	res, err = client.Do(req)
	if err != nil {
		util.Error(c, 500, "向GPT发送Run回答问题的请求失败", err)
		return
	}
	defer res.Body.Close()

	//

}

func GetThreads(c *gin.Context) {
	userid, err := ReadJwt(c)
	if err != nil {
		util.Error(c, 500, "JWT解析失败", nil)
		return
	}
	var user config.User
	err = config.DB.Preload("GptThreads").First(&user, userid).Error
	if err != nil {
		util.DbQueryError(c, err, "找不到对应的用户")
		return
	}
}

func ModifyThreads(c *gin.Context) {
	userid, err := ReadJwt(c)
	if err != nil {
		util.Error(c, 500, "JWT解析失败", err)
		return
	}
	action, threadId := c.PostForm("action"), c.PostForm("session_id")
	var delthread config.GptThread
	if err := config.DB.First(&delthread, "user_id = ? AND thread_id = ?", userid, threadId).Error; err != nil {
		util.DbQueryError(c, err, "找不到对应的会话")
		return
	}
	if action == "modify" {
		threadName := c.PostForm("thread_name")
		config.DB.Model(&delthread).Update("thread_name", threadName)
		util.Info(c, 200, "会话名更新成功", nil)
	} else if action == "delete" {
		config.DB.Delete(&delthread)
		util.Info(c, 200, "会话删除成功", nil)
	} else {
		util.Error(c, 400, "不支持的action或未指定action", nil)
		return
	}
}

func gptRequest(url, reqBody string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(reqBody)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer"+config.Config.GptToken)
	req.Header.Add("OpenAI-Beta", "assistants=v1")
	return req, nil
}
