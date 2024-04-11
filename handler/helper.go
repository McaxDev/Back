package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	var reqbody Reqres
	if err := c.BindJSON(&reqbody); err != nil {
		util.Error(c, 400, "无法解析你发送的请求", err)
		return
	}
	question, model, thread := reqbody.Info.Content, reqbody.GptModel, reqbody.Pid
	jwt, exist := c.Get("userID")
	if !exist {
		util.Error(c, 500, "读取用户信息失败", nil)
		return
	}
	userID := jwt.(int)
	var session config.GptThread
	if thread != "" {
		err := config.DB.Where("thread = ? AND user_id = ?", thread, userID).First(&session).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				util.Error(c, 400, "找不到对应的会话", err)
			}
			util.Error(c, 500, "查询会话失败", err)
			return
		}
	}
	var user config.User
	err := config.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			util.Error(c, 400, "找不到对应的用户", err)
		}
		util.Error(c, 500, "查询对应的用户失败", err)
		return
	}
	if (model == "4" && user.BlueCoin < 2) || user.WhiteCoin < 2 {
		util.Error(c, 400, "你的余额不足，无法提问", nil)
		return
	}
	reqJson := fmt.Sprintf(`{"role": "user", "content": "%s"}`, question)
	url := fmt.Sprintf("https://api.openai.com/v1/threads/%s/messages", thread)
	if thread == "" {
		url = "https://api.openai.com/v1/threads"
		reqJson = fmt.Sprintf(`{"messages": [%s]}`, reqJson)
	}
	payload := bytes.NewReader([]byte(reqJson))
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		util.Error(c, 500, "向GPT发送的请求创建失败", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.Config.GptToken)
	req.Header.Add("OpenAI-Beta", "assistants=v1")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		util.Error(c, 500, "向GPT发送请求失败", err)
		return
	}
	defer res.Body.Close()
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
	if temp, ok := data["id"].(string); !ok {
		util.Error(c, 500, "对thread id类型断言失败", nil)
		return
	} else {
		thread = temp
	}
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
