package serializer

import "todo_list/model"

// TaskSerializer 用于任务的序列化
type AdminTaskSerializer struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Status    int    `json:"status"`
	Uid       uint   `json:"uid"`
	UserName  string `json:"user_name"`
	Content   string `json:"content"`
	Priority  int    `json:"priority"`
	IsNotify  int    `json:"is_notify"`
	NotifyWay int    `json:"notify_way"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	IsChecked int    `json:"is_checked"`
}

// BuildTask 构建任务的序列化数据
func AdminBuildTask(task model.Task) AdminTaskSerializer {
	return AdminTaskSerializer{
		ID:        task.ID,
		Title:     task.Title,
		Status:    task.Status,
		Uid:       task.Uid,
		Content:   task.Content,
		UserName:  task.User.UserName, // 从关联的 User 取用户名
		Priority:  task.Priority,
		IsNotify:  task.IsNotify,
		NotifyWay: task.NotifyWay,
		StartTime: task.StartTime,
		EndTime:   task.EndTime,
		IsChecked: task.IsChecked,
	}
}
