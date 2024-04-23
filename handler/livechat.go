package handler

import (
	"context"
	"net/http"

	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 通过WebSocket协议将消息转发到网站的handler
func LiveChat(c *gin.Context) {

	// 从JWT里读取用户身份
	username := "未知用户"
	if user, err := BindJwt(c); err == nil {
		username = user.Username
	}

	// 建立websocket客户端
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		util.Error(c, 500, "WebSocket连接建立失败", err)
		return
	}
	defer ws.Close()

	// 从Redis里持续监听新的消息
	go func() {
		sub := co.RDB.Subscribe(context.Background(), "mainchat")
		defer sub.Close()
		for msg := range sub.Channel() {
			if err := ws.WriteMessage(1, []byte(msg.Payload)); err != nil {
				continue
			}
		}
	}()

	// 从客户端持续监听消息并存储到Redis
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			continue
		}
		msg = append([]byte(username), msg...)
		err = co.RDB.Publish(context.Background(), "webchat", msg).Err()
		if err != nil {
			continue
		}
	}
}

// 创建一个WebSocket连接创建器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
