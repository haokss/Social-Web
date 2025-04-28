package serializer

import (
	"todo_list/model"
)

type Task struct {
	ID         uint   `json:"id" example:"1"`
	Title      string `json:"title"`
	Status     int    `json:"code"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
}

// 返回单个Task
func BuildTask(item model.Task) Task {
	return Task{
		ID:         item.ID,
		Title:      item.Title,
		Status:     item.Status,
		Content:    item.Content,
		CreateTime: item.CreatedAt.Unix(),
		StartTime:  item.StartTime,
		EndTime:    item.EndTime,
	}
}

// 构建Task返回值切片
func BuildTasks(items []model.Task) []Task {
	var tasks []Task
	for _, item := range items {
		tasks = append(tasks, BuildTask(item))
	}
	return tasks
}
