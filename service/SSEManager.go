package service

import (
	"time"
	"todo_list/model"
)

type NotifyTimingTaskEvent struct {
	TaskID  uint   `json:"task_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Time    int64  `json:"time"` // 通知时间戳
}

// 全局通知通道
var NotificationChan = make(chan NotifyTimingTaskEvent)

// sendNotification 发送任务提醒通知
func sendNotification(task model.TimingTask) {
	// 构造通知事件
	event := NotifyTimingTaskEvent{
		TaskID:  task.ID,
		Title:   task.Title,
		Content: task.Content,
		Time:    time.Now().Unix(),
	}

	// 将事件发送到通知通道
	NotificationChan <- event
}
