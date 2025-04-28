package service

import (
	"time"
	"todo_list/model"
	"todo_list/serializer"
)

type CreateTaskService struct {
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
	Status  int    `json:"code" form:"code"` // 0未作，1已完成
}

type UpdateTaskService struct {
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
	Status  int    `json:"code" form:"code"` // 0未作，1已完成
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

// 创建一条备忘录
func (service *CreateTaskService) Create(id uint) serializer.Response {
	var user model.User
	model.DB.First(&user, id)
	code := 200
	task := model.Task{
		User:      user,
		Uid:       user.ID,
		Title:     service.Title,
		Status:    0,
		Content:   service.Content,
		StartTime: time.Now().Unix(),
		EndTime:   0,
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

// 展示一条备忘录
func (service *ShowTaskService) Show(uid uint, tid string) serializer.Response {
	var task model.Task
	code := 200
	err := model.DB.First(&task, tid).Error
	if err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg:    "search error!",
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildTask(task),
		Msg:    "search success!",
	}
}

// 返回所有的备忘录
func (service *ShowTaskAllService) ShowAll(uid uint) serializer.Response {
	var tasks []model.Task
	// 传入的pagesize为0，代表一次获取所有
	count := 0
	if service.PageNum == 0 {
		service.PageSize = 10
	}
	model.DB.Model(&model.Task{}).
		Preload("User").Where("uid=?", uid).
		Count(&count).
		Limit(service.PageSize).
		Offset((service.PageNum - 1) * service.PageSize).
		Find(&tasks)

	return serializer.BuildListResponse(serializer.BuildTasks(tasks), uint(count))
}

// 更新一条备忘录
func (service *UpdateTaskService) Update(uid uint, tid string) serializer.Response {
	var task model.Task
	code := 200
	err := model.DB.First(&task, tid).Error
	if err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg:    "search error!",
		}
	}
	task.Content = service.Content
	task.Title = service.Title
	task.Status = service.Status
	model.DB.Save(&task)
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildTask(task),
		Msg:    "update success!",
	}
}

// 模糊查询
func (service *SearchTaskService) Search(uid uint) serializer.Response {
	var tasks []model.Task
	count := 0

	model.DB.Model(&model.Task{}).
		Preload("User").
		Where("title LIKE ? OR content LIKE ?", "%"+service.SearchInfo+"%", "%"+service.SearchInfo+"%").
		Count(&count).
		Limit(service.PageSize).
		Offset((service.PageNum - 1) * service.PageSize).
		Find(&tasks)

	return serializer.Response{
		Status: 200,
		Data:   serializer.BuildTasks(tasks),
	}
}

// 删除
func (service *DeleteTaskService) Delete(tid string) serializer.Response {
	var task model.Task
	model.DB.First(&task, tid)
	err := model.DB.Delete(&task).Error
	if err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "Delete Task Failed!",
		}
	}
	return serializer.Response{
		Status: 200,
		Msg:    "Delete Task Success!",
	}
}
