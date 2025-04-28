package api

import (
	"todo_list/service"

	"github.com/gin-gonic/gin"
)

// Web推送示例（使用SSE）
func PushUpdates(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	for {
		select {
		case event := <-service.NotificationChan:
			c.SSEvent("message", event) // 发送事件到客户端
			c.Writer.Flush()
		case <-c.Done():
			return // 客户端断开连接
		}
	}
}
