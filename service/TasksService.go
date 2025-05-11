package service

import (
	"todo_list/cache"
	"todo_list/model"
	sse "todo_list/package/SSE"
	"todo_list/package/utils"
	"todo_list/serializer"

	"github.com/gin-gonic/gin"
)

type CreateTaskService struct {
	Title     string `json:"title" form:"title"`
	Content   string `json:"content" form:"content"`
	Status    int    `json:"code" form:"code"` // 0未作，1已完成
	Priority  int    `json:"priority" form:"priority"`
	IsNotify  int    `json:"IsNotify" form:"IsNotify"`
	NotifyWay int    `json:"NotifyWay" form:"NotifyWay"`
	StartTime int64  `json:"start_time" form:"start_time"`
	EndTime   int64  `json:"end_time" form:"end_time"`
}

type UpdateTaskService struct {
	Title     string `json:"title" form:"title"`
	Content   string `json:"content" form:"content"`
	Status    int    `json:"code" form:"code"` // 0未作，1已完成
	Priority  int    `json:"priority" form:"priority"`
	IsNotify  int    `json:"IsNotify" form:"IsNotify"`
	NotifyWay int    `json:"NotifyWay" form:"NotifyWay"`
	StartTime int64  `json:"start_time" form:"start_time"`
	EndTime   int64  `json:"end_time" form:"end_time"`
}

type SearchTaskService struct {
	SearchInfo string `json:"search_info" form:"search_info"`
	PageNum    int    `json:"page_num" form:"page_num"`
	PageSize   int    `json:"page_size" from:"page_size"`
}

type ShowTaskService struct {
}

type DeleteTaskService struct {
}

type ShowTaskAllService struct {
	PageNum  int `json:"page_num" form:"page_num"`
	PageSize int `json:"page_size" from:"page_size"`
}

// 创建一条活动
func (service *CreateTaskService) Create(c *gin.Context, id uint) serializer.Response {
	var user model.User
	model.DB.First(&user, id)
	task := model.Task{
		User:      user,
		Uid:       user.ID,
		Title:     service.Title,
		Status:    0,
		Content:   service.Content,
		Priority:  service.Priority,
		StartTime: service.StartTime,
		EndTime:   service.EndTime,
		IsNotify:  service.IsNotify,
		NotifyWay: service.NotifyWay,
	}
	if err := model.DB.Create(&task).Error; err != nil {
		return serializer.Response{Status: 500, Msg: "create failed!"}
	}
	cache.SetTask(task)

	broker := c.MustGet("sseBroker").(*sse.Broker)

	// 通知管理员审核
	broker.Notify(sse.Message{
		Event:     "instant_notification",
		Data:      map[string]interface{}{"title": "您有新的活动待审核：" + service.Title, "content": service.Content},
		TargetIDs: []uint{utils.AdminUid},
	})
	return serializer.Response{Status: 200, Msg: "create success!"}
}

// 展示一条活动
func (service *ShowTaskService) Show(uid uint, tid string) serializer.Response {
	taskID := utils.ParseID(tid)
	if task, ok := cache.GetTask(taskID); ok && task.Uid == uid {
		return serializer.Response{Status: 200, Data: serializer.BuildTask(task), Msg: "search success!"}
	}
	return serializer.Response{Status: 404, Msg: "task not found!"}
}

// 返回所有的活动
func (service *ShowTaskAllService) ShowAll(uid uint) serializer.Response {
	tasks := cache.GetUserTasks(uid)

	// 过滤掉审核未通过的任务
	filtered := make([]model.Task, 0, len(tasks))
	for _, task := range tasks {
		if task.IsChecked != 2 {
			filtered = append(filtered, task)
		}
	}

	count := len(filtered)
	return serializer.BuildListResponse(serializer.BuildTasks(filtered), uint(count))
}

// 更新一条活动
func (service *UpdateTaskService) Update(uid uint, tid string) serializer.Response {
	taskID := utils.ParseID(tid)
	var task model.Task
	if err := model.DB.First(&task, taskID).Error; err != nil {
		return serializer.Response{Status: 500, Msg: "search error!"}
	}
	task.Title = service.Title
	task.Content = service.Content
	task.Status = service.Status
	task.Priority = service.Priority
	task.StartTime = service.StartTime
	task.EndTime = service.EndTime
	task.IsNotify = service.IsNotify
	task.NotifyWay = service.NotifyWay
	model.DB.Save(&task)
	cache.SetTask(task)
	return serializer.Response{Status: 200, Data: serializer.BuildTask(task), Msg: "update success!"}
}

// 模糊查询
func (service *SearchTaskService) Search(uid uint) serializer.Response {
	tasks := cache.GetUserTasks(uid)
	var result []model.Task
	for _, t := range tasks {
		if utils.Contains(t.Title, service.SearchInfo) || utils.Contains(t.Content, service.SearchInfo) {
			result = append(result, t)
		}
	}
	return serializer.Response{Status: 200, Data: serializer.BuildTasks(result)}
}

// 删除活动
func (service *DeleteTaskService) Delete(tid string) serializer.Response {
	taskID := utils.ParseID(tid)
	if err := model.DB.Unscoped().Delete(&model.Task{}, taskID).Error; err != nil {
		return serializer.Response{Status: 500, Msg: "Delete Task Failed!"}
	}
	cache.DeleteTask(taskID)
	return serializer.Response{Status: 200, Msg: "Delete Task Success!"}
}
