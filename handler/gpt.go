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

// 代表GPT的会话格式的结构体
type GptSession struct {
	UserID    uint         `json:"userId"`
	Username  string       `json:"username"`
	SessionID string       `json:"sessionId"`
	Message   []GptMessage `json:"message"`
}

// 代表GPT的消息格式的结构体
type GptMessage struct {
	MessageID string `json:"messageId"`
	Role      string `json:"role"`
	Time      string `json:"time"`
	GptModel  string `json:"gptModel"`
	Content   string `json:"content"`
}

// 创建GPT连接客户端
var cli *ai.Client

// 向GPT提问的handler
func Gpt(c *gin.Context) {

	// 从请求体获得数据
	var req struct {
		ThreadID string `json:"sessionId"`
		GptModel string `json:"gptModel"`
		Message  string `json:"message"`
	}
	if err := util.BindReq(c, &req); err != nil {
		util.Error(c, 400, "你的请求体格式不正确", err)
		return
	}

	// 从JWT里获取用户ID
	user, err := BindJwt(c, "Balance", "Thread")
	if err != nil {
		util.Error(c, 500, "无法读取你的用户信息", err)
		return
	}

	// 根据用户的请求读取AssistantID
	asst := ai.RunRequest{AssistantID: co.Config.AsstID["GPT3.5"]}
	cost := -1
	if req.GptModel == "HELPER" {
		asst.AssistantID = co.Config.AsstID["HELPER"]
	} else if req.GptModel == "GPT4" {
		asst.AssistantID = co.Config.AsstID["GPT4"]
		cost = -2
	}

	// 检查用户的钱够不够并扣钱
	if err := user.Transact(cost); err != nil {
		util.Error(c, 400, "扣费失败", err)
		return
	}

	if req.ThreadID == "" { // 创建新的会话

		// 创建一个等待会话超时的上下文
		ctx, canc := util.Timeout(30)
		defer canc()

		// 创建会话
		thread, err := cli.CreateThread(ctx, ai.ThreadRequest{
			Messages: []ai.ThreadMessage{{
				Role:    ai.ThreadMessageRole("user"),
				Content: req.Message,
			}}})
		if err != nil {
			util.Error(c, 500, "会话创建失败", err)
			return
		}

		// 将请求里的会话ID修改为新创建的会话ID
		req.ThreadID = thread.ID

		// 将用户的会话信息存储到数据库
		if err := co.DB.Create(&co.GptThread{
			ThreadID:   thread.ID,
			ThreadName: time.Now().Format("2006-01-02 15:04:05"),
			UserID:     user.ID,
		}).Error; err != nil {
			util.DbQueryError(c, err, "无法将你的会话信息存储到数据库")
			return
		}

	} else { // 使用已有的会话

		// 检查用户是否拥有这个会话
		result := co.DB.First(&co.GptThread{}, "user_id = ? AND thread_id = ?", user.ID, req.ThreadID)
		if err := result.Error; err != nil {
			util.Error(c, 400, "你没有这个会话", err)
			return
		}

		// 创建一个等待会话超时的上下文
		ctx, canc := util.Timeout(30)
		defer canc()

		// 将用户的消息添加到会话里
		if _, err := cli.CreateMessage(ctx, req.ThreadID, ai.MessageRequest{
			Role:    "user",
			Content: req.Message,
		}); err != nil {
			util.Error(c, 500, "无法将你的消息添加到会话", err)
			return
		}

	}

	// 创建一个等待会话超时的上下文
	ctx, canc := util.Timeout(30)
	defer canc()

	// 生成回答
	run, err := cli.CreateRun(ctx, req.ThreadID, asst)
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
	util.Info(c, 200, "执行成功", GptSession{
		UserID:    user.ID,
		Username:  user.Username,
		SessionID: req.ThreadID,
		Message:   mes,
	})
}

// 修改会话名称或删除会话或列出所有会话或列出会话消息
func GptUtil(c *gin.Context) {

	// 从用户的请求体里获得查询字符串参数
	threadID, threadName, action := c.Query("session_id"), c.Query("session_name"), c.Query("action")

	// 从JWT里获取用户ID
	user, err := BindJwt(c, "Thread")
	if err != nil {
		util.Error(c, 500, "读取用户JWT信息失败", err)
		return
	}

	if threadID == "" { // 如果会话ID为空，就返回这个用户的所有会话

		// 将切片里的所有会话映射到结构体里
		var results []map[string]any
		for _, thread := range user.Thread {
			results = append(results, map[string]any{
				"session_id":   thread.ThreadID,
				"session_name": thread.ThreadName,
				"time":         thread.UpdatedAt.In(util.Loc).Format("2006/01/02 15:04"),
			})
		}
		util.Info(c, 200, "会话信息查找完成", results)

	} else { // 如果会话ID不为空，进行删除会话或修改会话名或列出消息

		// 检查用户是否拥有这个会话
		var tmp co.GptThread
		err = co.DB.First(&tmp, "thread_id = ? AND user_id = ?", threadID, user.ID).Error
		if err != nil {
			util.DbQueryError(c, err, "你没有这个会话")
			return
		}

		if action == "delete" { // 如果行为等于删除，删除会话
			result := co.DB.Delete(&tmp, "thread_id = ?", threadID)
			if err := result.Error; err != nil {
				util.Error(c, 500, "无法删除这个会话", err)
				return
			}
			util.Info(c, 200, "会话删除成功", nil)

		} else if threadName != "" { // 如果会话名称不为空，修改会话名称
			result := co.DB.Model(&tmp).Update("thread_name", threadName)
			if err := result.Error; err != nil {
				util.Error(c, 500, "无法修改这个会话的名称", err)
				return
			}
			util.Info(c, 200, "会话名称修改成功", nil)

		} else { // 否则，查看会话内容
			messages, err := FetchMessages(cli, threadID)
			if err != nil {
				util.Error(c, 500, "无法获取会话消息", err)
				return
			}
			util.Info(c, 200, "查询成功", messages)
		}
	}
}

// 对run对象进行轮询检测判断执行状态
func PollRunStatus(cli *ai.Client, threadID, runID string) (mes []GptMessage, err error) {
	ticker := time.NewTicker(1 * time.Second)
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
			case ai.RunStatusQueued, ai.RunStatusInProgress:
				continue
			default:
				return nil, fmt.Errorf("run failed with status: %s", run.Status)
			}
		}
	}
}

// 从openai的thread里提取内容到GptMessage切片
func FetchMessages(cli *ai.Client, threadID string) (destination []GptMessage, err error) {

	// 创建一个等待会话超时的上下文
	ctx, canc := util.Timeout(30)
	defer canc()

	messages, err := cli.ListMessage(ctx, threadID, nil, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("无法列出消息：%w", err)
	}
	for _, value := range messages.Messages {
		destination = append(destination, GptMessage{
			MessageID: value.ID,
			Role:      value.Role,
			Time:      time.Now().Format("2006-01-02 15:04:05"),
			GptModel:  DetectGptModel(util.Deref(value.AssistantID)),
			Content:   value.Content[0].Text.Value,
		})
	}
	return
}

// 将GPT的id翻译为gpt模型号的函数
func DetectGptModel(asstid string) string {
	for key, value := range co.Config.AsstID {
		if asstid == value {
			return key
		}
	}
	return ""
}
