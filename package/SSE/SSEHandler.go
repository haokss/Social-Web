package sse

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// SSEHandler SSE连接处理
func SSEHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户信息
		userID := c.MustGet("currentUserID").(uint)
		broker := c.MustGet("sseBroker").(*Broker)

		// 初始化客户端
		client := &Client{
			UserID:  userID,
			Message: make(chan Message, 10),
		}

		// 注册客户端
		broker.Register(client)
		defer broker.Unregister(userID, client)

		// 设置SSE响应头
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")

		// 保持连接
		c.Stream(func(w io.Writer) bool {
			select {
			case msg := <-client.Message:
				c.SSEvent(msg.Event, msg.Data)
			case <-c.Writer.CloseNotify():
				return false
			case <-time.After(30 * time.Second):
				c.SSEvent("ping", nil) // 发送心跳
			}
			return true
		})
	}
}
