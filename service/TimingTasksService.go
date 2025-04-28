package service

import (
	"time"
	"todo_list/model"
	"todo_list/serializer"
)

type CreateTimingTaskService struct {
	Title     string    `json:"title" form:"title"`
	Content   string    `json:"content" form:"content"`
	Priority  int       `json:"priority" form:"priority"`
	NotifyWay int       `json:"notify_way" form:"notify_way"`
	EarlyTime int       `json:"early_time" form:"early_time"`
	StartTime time.Time `json:"start_time" form:"start_time"`
	EndTime   time.Time `json:"end_time" form:"start_time"`
	Type      int       `json:"type" form:"type"`
}

type UpdateTimingTaskService struct {
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
	Status  int    `json:"code" form:"code"` // 0未作，1已完成
}

type SearchTimingTaskService struct {
	SearchInfo string `json:"search_info" form:"search_info"`
	PageNum    int    `json:"page_num" form:"page_num"`
	PageSize   int    `json:"page_size" from:"page_size"`
}

type ShowTimingTaskService struct {
}

type DeleteTimingTaskService struct {
}

type ShowTimingTaskAllService struct {
	PageNum  int `json:"page_num" form:"page_num"`
	PageSize int `json:"page_size" from:"page_size"`
}

// 任务执行状态
const (
	TaskStatusPending   = 0 // 未开始
	TaskStatusRunning   = 1 // 进行中
	TaskStatusCompleted = 2 // 已完成
	TaskStatusExpired   = 3 // 已过期
)

// 创建一条备忘录
func (service *CreateTimingTaskService) Create(id uint) serializer.Response {
	var user model.User
	model.DB.First(&user, id)
	code := 200
	task := model.TimingTask{
		User:      user,
		Uid:       user.ID,
		Title:     service.Title,
		Status:    0,
		Content:   service.Content,
		Priority:  service.Priority,
		Type:      service.Type,
		NotifyWay: service.NotifyWay,
		EarlyTime: service.EarlyTime,
		StartTime: service.StartTime,
		EndTime:   service.EndTime,
	}
	err := model.DB.Create(&task).Error
	if err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg:    "create failed!",
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    "create success!",
	}
}

// // 展示一条定时任务
// func (service *ShowTimingTaskService) Show(uid uint, tid string) serializer.Response {
// 	// var task model.TimingTask
// 	// code := 200
// 	// err := model.DB.First(&task, tid).Error
// 	// if err != nil {
// 	// 	code = 500
// 	// 	return serializer.Response{
// 	// 		Status: code,
// 	// 		Msg:    "search error!",
// 	// 	}
// 	// }
// 	// return serializer.Response{
// 	// 	Status: code,
// 	// 	Data:   serializer.BuildTask(task),
// 	// 	Msg:    "search success!",
// 	// }
// }

// 返回所有定时任务
func (service *ShowTimingTaskAllService) ShowAll(uid uint) serializer.Response {
	var tasks []model.TimingTask
	// 传入的pagesize为0，代表一次获取所有
	count := 0
	if service.PageNum == 0 {
		service.PageSize = 10
	}

	model.DB.Model(&model.TimingTask{}).
		Preload("User").Where("uid=?", uid).
		Count(&count).
		Limit(service.PageSize).
		Offset((service.PageNum - 1) * service.PageSize).
		Find(&tasks)

	return serializer.Response{
		Status: 200,
		Data: serializer.DataList{
			Item:  serializer.BuildTimingTasks(tasks),
			Total: uint(count),
		},
		Msg: "get success!",
	}
}

// // 更新一条备忘录
// func (service *UpdateTimingTaskService) Update(uid uint, tid string) serializer.Response {
// 	// var task model.TimingTask
// 	// code := 200
// 	// err := model.DB.First(&task, tid).Error
// 	// if err != nil {
// 	// 	code = 500
// 	// 	return serializer.Response{
// 	// 		Status: code,
// 	// 		Msg:    "search error!",
// 	// 	}
// 	// }
// 	// task.Content = service.Content
// 	// task.Title = service.Title
// 	// task.Status = service.Status
// 	// model.DB.Save(&task)
// 	// return serializer.Response{
// 	// 	Status: code,
// 	// 	Data:   serializer.BuildTask(task),
// 	// 	Msg:    "update success!",
// 	// }
// }

// // 模糊查询
// func (service *SearchTimingTaskService) Search(uid uint) serializer.Response {
// 	// var tasks []model.TimingTask
// 	// count := 0

// 	// model.DB.Model(&model.Task{}).
// 	// 	Preload("User").
// 	// 	Where("title LIKE ? OR content LIKE ?", "%"+service.SearchInfo+"%", "%"+service.SearchInfo+"%").
// 	// 	Count(&count).
// 	// 	Limit(service.PageSize).
// 	// 	Offset((service.PageNum - 1) * service.PageSize).
// 	// 	Find(&tasks)

// 	// return serializer.Response{
// 	// 	Status: 200,
// 	// 	Data:   serializer.BuildTasks(tasks),
// 	// }
// }

// // 删除
// func (service *DeleteTimingTaskService) Delete(tid string) serializer.Response {
// 	// var task model.TimingTask
// 	// model.DB.First(&task, tid)
// 	// err := model.DB.Delete(&task).Error
// 	// if err != nil {
// 	// 	return serializer.Response{
// 	// 		Status: 500,
// 	// 		Msg:    "Delete Task Failed!",
// 	// 	}
// 	// }
// 	// return serializer.Response{
// 	// 	Status: 200,
// 	// 	Msg:    "Delete Task Success!",
// 	// }
// }

// 定时扫描未完成的任务
// func StartTaskScheduler() {
// 	ticker := time.NewTicker(30 * time.Second)
// 	defer ticker.Stop()
// 	for range ticker.C {
// 		now := time.Now().UTC()
// 		var tasks []model.TimingTask
// 		model.DB.Where(
// 			"status IN (?) AND DATE_SUB(start_time, INTERVAL early_time SECOND) <= ? AND end_time >= ?",
// 			[]int{TaskStatusPending, TaskStatusRunning},
// 			now,
// 			now,
// 		).Find(&tasks)
// 		for _, task := range tasks {
// 			sendNotification(task) // 调用提醒服务
// 			fmt.Println(task.Content)
// 		}

// 		// 标记过期任务
// 		model.DB.Model(&model.TimingTask{}).
// 			Where("end_time < ? AND status != ?", now, TaskStatusCompleted).
// 			Update("status", TaskStatusExpired)
// 	}
// }

// Web推送示例（使用SSE）
// func (s *TaskService) PushUpdates(c *gin.Context) {
// 	c.Header("Content-Type", "text/event-stream")
// 	for {
// 		// 监听任务状态变更事件
// 		select {
// 		case event := <-notificationChan:
// 			c.SSEvent("message", event)
// 			c.Writer.Flush()
// 		case <-c.Done():
// 			return
// 		}
// 	}
// }
