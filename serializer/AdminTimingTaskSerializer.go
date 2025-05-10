package serializer

import (
	"time"
	"todo_list/model"
)

type AdminTimingTaskSerializer struct {
	ID        uint      `json:"id"`
	Uid       uint      `json:"uid"`
	Username  string    `json:"user_name"`
	Title     string    `json:"title"`
	Status    int       `json:"status"`
	Content   string    `json:"content"`
	Priority  int       `json:"priority"`
	Type      int       `json:"type"`
	NotifyWay int       `json:"notify_way"`
	EarlyTime int       `json:"early_time"`
	IsChecked int       `json:"is_checked"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func AdminBuildTimingTask(task model.TimingTask) AdminTimingTaskSerializer {
	return AdminTimingTaskSerializer{
		ID:        task.ID,
		Uid:       task.Uid,
		Username:  task.User.UserName,
		Title:     task.Title,
		Status:    task.Status,
		Content:   task.Content,
		Priority:  task.Priority,
		Type:      task.Type,
		NotifyWay: task.NotifyWay,
		EarlyTime: task.EarlyTime,
		IsChecked: task.IsChecked,
		StartTime: task.StartTime,
		EndTime:   task.EndTime,
	}
}
