package handler

import (
	"context"
	"fmt"
	"time"

	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	ai "github.com/sashabaranov/go-openai"
)

// 创建GPT连接客户端
// var cli = ai.NewClient(co.Config.GptToken)
var loggingClient = util.NewLoggingClient() // 使用自定义的日志记录客户端
var cli = ai.NewClientWithHTTPClient(co.Config.GptToken, loggingClient)

// 向GPT提问的handler
func Gpt(c *gin.Context) {

	// 从请求体获得数据
	var req struct {
		ThreadID string `json:"thread_id"`
		GptModel string `json:"gpt_model"`
		Message  string `json:"message"`
	}
	if err := util.BindReq(c, &req); err != nil {
		util.Error(c, 400, "你的请求体格式不正确", err)
		return
	}

	// 从JWT里获取用户ID
	userID, err := ReadJwt(c)
	if err != nil {
		util.Error(c, 500, "无法读取你的用户信息", err)
		return
	}

	if req.ThreadID == "" { // 创建新的会话

		// 创建会话
		thread, err := cli.CreateThread(util.Timeout(30), ai.ThreadRequest{
			Messages: []ai.ThreadMessage{{
				Role:    ai.ThreadMessageRole("user"),
				Content: req.Message,
			}}})
		if err != nil {
			util.Error(c, 500, "会话创建失败", err)
			return
		}
		req.ThreadID = thread.ID

		// 将用户的会话信息存储到数据库
		if err := co.DB.Create(co.GptThread{
			ThreadID:   thread.ID,
			ThreadName: time.Now().Format("2006-01-02 15:04:05"),
			UserID:     userID,
		}).Error; err != nil {
			util.Error(c, 500, "无法将你的会话信息存储到数据库", err)
			return
		}

	} else { // 使用已有的会话

		// 检查用户是否拥有这个会话
		err := co.DB.First(co.GptThread{}, "user_id = ? AND thread_id = ?", userID, req.ThreadID).Error
		if err != nil {
			util.Error(c, 400, "你没有这个会话", err)
			return
		}

		// 将用户的消息添加到会话里
		if _, err = cli.CreateMessage(util.Timeout(30), req.ThreadID, ai.MessageRequest{
			Role:    "user",
			Content: req.Message,
		}); err != nil {
			util.Error(c, 500, "无法将你的消息添加到会话", err)
			return
		}

	}

	// 根据用户的请求读取AssistantID
	var asst_id string
	switch req.GptModel {
	case "GPT3.5":
		asst_id = co.Config.Gpt3
	case "GPT4":
		asst_id = co.Config.Gpt4
	case "HELPER":
		asst_id = co.Config.Axo
	}

	// 生成回答
	run, err := cli.CreateRun(util.Timeout(30), req.ThreadID, ai.RunRequest{
		AssistantID: asst_id,
	})
	if err != nil {
		util.Error(c, 500, "无法生成回答", err)
		return
	}

	// 对run对象进行轮询检测，判断执行状态
	mes, err := PollRunStatus(cli, req.ThreadID, run.ID)
	if err != nil {
		util.Error(c, 500, "获取消息失败", err)
		return
	}

	// 将回答返回给用户
	util.Info(c, 200, "执行成功", mes)
}

func PollRunStatus(cli *ai.Client, threadID, runID string) ([]ai.Message, error) {
	ticker := time.NewTicker(5 * time.Second) // Check every 5 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			run, err := cli.RetrieveRun(context.Background(), threadID, runID)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve run: %w", err)
			}

			switch run.Status {
			case ai.RunStatusCompleted:
				return FetchMessages(cli, threadID)
			case ai.RunStatusFailed, ai.RunStatusCancelled:
				return nil, fmt.Errorf("run failed or was cancelled")
			// Continue polling if still in progress or queued
			case ai.RunStatusQueued, ai.RunStatusInProgress:
				continue
			default:
				return nil, fmt.Errorf("unexpected run status: %s", run.Status)
			}
		}
	}
}

func FetchMessages(cli *ai.Client, threadID string) ([]ai.Message, error) {
	messages, err := cli.ListMessage(context.Background(), threadID, nil, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	return messages.Messages, nil
}
