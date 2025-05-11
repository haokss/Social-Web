package service

import (
	"time"
	"todo_list/cache"
	"todo_list/model"
	sse "todo_list/package/SSE"
	"todo_list/package/utils"
	"todo_list/serializer"

	"github.com/gin-gonic/gin"
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

type DeleteTimingTasksService struct {
	IDs []uint `json:"ids" form:"ids"` // 任务ID数组
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

// 创建一条定时任务
func (service *CreateTimingTaskService) Create(c *gin.Context, id uint) serializer.Response {
	var user model.User
	model.DB.First(&user, id)
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
	if err := model.DB.Create(&task).Error; err != nil {
		return serializer.Response{Status: 500, Msg: "create failed!"}
	}
	cache.SetTimingTask(task) // 加入缓存

	// 创建定时任务通知
	// 计算提醒时间
	notifyTime := service.StartTime.Add(-time.Duration(service.EarlyTime) * time.Minute)
	broker := c.MustGet("sseBroker").(*sse.Broker)
	broker.ScheduleNotify(sse.Message{
		Event:     "new_year_notification",
		Data:      map[string]interface{}{"title": service.Title, "content": service.Content},
		TargetIDs: []uint{user.ID},
	}, notifyTime)

	// 通知管理员审核
	broker.Notify(sse.Message{
		Event:     "instant_notification",
		Data:      map[string]interface{}{"title": "您有新的提醒待审核：" + service.Title, "content": service.Content},
		TargetIDs: []uint{utils.AdminUid},
	})

	return serializer.Response{Status: 200, Msg: "create success!"}
}

// 返回所有定时任务
func (service *ShowTimingTaskAllService) ShowAll(c *gin.Context, uid uint) serializer.Response {
	var tasks []model.TimingTask
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

	// 更新缓存
	for _, task := range tasks {
		cache.SetTimingTask(task)
	}

	return serializer.Response{
		Status: 200,
		Data: serializer.DataList{
			Item:  serializer.BuildTimingTasks(tasks),
			Total: uint(count),
		},
		Msg: "get success!",
	}
}

// 更新一条备忘录
func (service *UpdateTimingTaskService) Update(uid uint, tid string) serializer.Response {
	var task model.TimingTask
	if err := model.DB.Where("uid = ? AND id = ?", uid, tid).First(&task).Error; err != nil {
		return serializer.Response{Status: 500, Msg: "task not found!"}
	}

	task.Title = service.Title
	task.Content = service.Content
	task.Status = service.Status

	if err := model.DB.Save(&task).Error; err != nil {
		return serializer.Response{Status: 500, Msg: "update failed!"}
	}
	cache.SetTimingTask(task)
	return serializer.Response{
		Status: 200,
		Data:   serializer.BuildTimingTask(task),
		Msg:    "update success!",
	}
}

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

// 删除定时任务
func (service *DeleteTimingTasksService) Delete(uid uint) serializer.Response {
	if len(service.IDs) == 0 {
		return serializer.Response{Status: 400, Msg: "请选择要删除的任务"}
	}
	if err := model.DB.Where("uid = ? AND id IN (?)", uid, service.IDs).Delete(&model.TimingTask{}).Error; err != nil {
		return serializer.Response{Status: 500, Msg: "批量删除失败"}
	}
	for _, id := range service.IDs {
		cache.DeleteTimingTask(id)
	}
	return serializer.Response{Status: 200, Msg: "批量删除成功"}
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
