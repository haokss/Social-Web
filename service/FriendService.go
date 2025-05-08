package service

import (
	"todo_list/model"
	"todo_list/serializer"
)

type CreateFriendService struct {
	Name     string `json:"name" form:"name"`
	Phone    string `json:"phone" form:"phone"`
	Birthday string `json:"birthday" form:"birthday"`
	Hobby    string `json:"hobby" form:"hobby"`
}

type ListFriendService struct {
	Name     string `form:"name"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type UpdateFriendService struct {
	Name     string `json:"name" form:"name"`
	Phone    string `json:"phone" form:"phone"`
	Birthday string `json:"birthday" form:"birthday"`
	Hobby    string `json:"hobby" form:"hobby"`
}

type BatchDeleteFriendService struct {
	IDs []uint `json:"ids"`
}

// 新建朋友
func (service *CreateFriendService) Create(uid uint) serializer.Response {
	friend := model.Friend{
		Uid:      uid,
		Name:     service.Name,
		Phone:    service.Phone,
		Birthday: service.Birthday,
		Hobby:    service.Hobby,
	}
	if err := model.DB.Create(&friend).Error; err != nil {
		return serializer.Response{Status: 500, Msg: "创建失败"}
	}
	return serializer.Response{Status: 200, Msg: "创建成功"}
}

// 查询朋友列表
func (service *ListFriendService) List(uid uint) serializer.Response {
	var friends []model.Friend
	var total int64

	if service.Page <= 0 {
		service.Page = 1
	}
	if service.PageSize <= 0 {
		service.PageSize = 10
	}

	tx := model.DB.Model(&model.Friend{}).Where("uid = ?", uid)
	if service.Name != "" {
		tx = tx.Where("name LIKE ?", "%"+service.Name+"%")
	}
	tx.Count(&total)

	err := tx.Order("created_at desc").
		Offset((service.Page - 1) * service.PageSize).
		Limit(service.PageSize).
		Find(&friends).Error

	if err != nil {
		return serializer.Response{Status: 500, Msg: "获取朋友列表失败"}
	}

	return serializer.Response{
		Status: 200,
		Msg:    "获取朋友列表成功",
		Data: map[string]interface{}{
			"list":      friends,
			"total":     total,
			"page":      service.Page,
			"page_size": service.PageSize,
		},
	}
}

// 更新朋友
func (service *UpdateFriendService) Update(idStr string) serializer.Response {
	var friend model.Friend
	if err := model.DB.First(&friend, idStr).Error; err != nil {
		return serializer.Response{Status: 404, Msg: "朋友未找到"}
	}

	friend.Name = service.Name
	friend.Phone = service.Phone
	friend.Birthday = service.Birthday
	friend.Hobby = service.Hobby

	if err := model.DB.Save(&friend).Error; err != nil {
		return serializer.Response{Status: 500, Msg: "更新失败"}
	}
	return serializer.Response{Status: 200, Msg: "更新成功"}
}

// 批量删除
func (service *BatchDeleteFriendService) BatchDelete() serializer.Response {
	if len(service.IDs) == 0 {
		return serializer.Response{Status: 400, Msg: "请传入要删除的ID列表"}
	}
	if err := model.DB.Delete(&model.Friend{}, service.IDs).Error; err != nil {
		return serializer.Response{Status: 500, Msg: "删除失败"}
	}
	return serializer.Response{Status: 200, Msg: "批量删除成功"}
}
