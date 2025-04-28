package serializer

import (
	"time"
	"todo_list/model"
)

type TimingTask struct {
	ID         uint      `json:"id" example:"1"`
	Title      string    `json:"title"`
	Status     int       `json:"code"` // 0未完成，1 已完成
	Content    string    `json:"content"`
	CreateTime time.Time `json:"create_time"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`

	Priority  int `json:"priority"`   // 任务优先级：0~5
	Type      int `json:"type"`       // 任务类型
	NotifyWay int `json:"notify_way"` // 新增：提醒方式（0站内信/1邮件/2短信）
	EarlyTime int `json:"early_time"` // 提前的时间，秒
}

type DataList struct {
	Item  interface{} `json:"item"`
	Total uint        `json:"total"`
}

// 返回单个TimingTask
func BuildTimingTask(item model.TimingTask) TimingTask {
	return TimingTask{
		ID:         item.ID,
		Title:      item.Title,
		Status:     item.Status,
		Content:    item.Content,
		CreateTime: item.CreatedAt,
		StartTime:  item.StartTime,
		EndTime:    item.EndTime,
		Priority:   item.Priority,
		Type:       item.Type,
		NotifyWay:  item.NotifyWay,
		EarlyTime:  item.EarlyTime,
	}
}

// 构建Task返回值切片
func BuildTimingTasks(items []model.TimingTask) []TimingTask {
	var tasks []TimingTask
	for _, item := range items {
		tasks = append(tasks, BuildTimingTask(item))
	}
	return tasks
}
