package service

import (
	"todo_list/model"
	"todo_list/serializer"
)

type CreateClassmateService struct {
	Name      string `json:"name" form:"name"`
	Phone     string `json:"phone" form:"phone"`
	School    string `json:"school" form:"school"`
	ClassName string `json:"className" form:"className"`
	Stage     string `json:"stage" form:"stage"`
}
type ListClassmateService struct {
	Name     string `form:"name"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type UpdateClassmateService struct {
	Name      string `json:"name" form:"name"`
	Phone     string `json:"phone" form:"phone"`
	School    string `json:"school" form:"school"`
	ClassName string `json:"className" form:"className"`
	Stage     string `json:"stage" form:"stage"`
}

type BatchDeleteClassmateService struct {
	IDs []uint `json:"ids" form:"ids"`
}

func (service *CreateClassmateService) Create(uid uint) serializer.Response {
	classmate := model.Classmate{
		Uid:       uid,
		Name:      service.Name,
		Phone:     service.Phone,
		School:    service.School,
		ClassName: service.ClassName,
		Stage:     service.Stage,
	}
	if err := model.DB.Create(&classmate).Error; err != nil {
		return serializer.Response{Status: 500, Msg: "创建失败"}
	}
	return serializer.Response{Status: 200, Msg: "创建成功"}
}

func (service *ListClassmateService) List(uid uint) serializer.Response {
	var classmates []model.Classmate
	var total int64

	if service.Page <= 0 {
		service.Page = 1
	}
	if service.PageSize <= 0 {
		service.PageSize = 10
	}

	tx := model.DB.Model(&model.Classmate{}).Where("uid = ?", uid)
	if service.Name != "" {
		tx = tx.Where("name LIKE ?", "%"+service.Name+"%")
	}

	tx.Count(&total)

	err := tx.Order("created_at desc").
		Offset((service.Page - 1) * service.PageSize).
		Limit(service.PageSize).
		Find(&classmates).Error

	if err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "获取同学列表失败",
		}
	}

	return serializer.Response{
		Status: 200,
		Msg:    "获取同学列表成功",
		Data: map[string]interface{}{
			"list":      classmates,
			"total":     total,
			"page":      service.Page,
			"page_size": service.PageSize,
		},
	}
}

func (service *UpdateClassmateService) Update(idStr string) serializer.Response {
	var classmate model.Classmate
	if err := model.DB.First(&classmate, idStr).Error; err != nil {
		return serializer.Response{Status: 404, Msg: "同学未找到"}
	}
	classmate.Name = service.Name
	classmate.Phone = service.Phone
	classmate.School = service.School
	classmate.ClassName = service.ClassName
	classmate.Stage = service.Stage

	if err := model.DB.Save(&classmate).Error; err != nil {
		return serializer.Response{Status: 500, Msg: "更新失败"}
	}
	return serializer.Response{Status: 200, Msg: "更新成功"}
}

// 批量删除
func (service *BatchDeleteClassmateService) BatchDelete() serializer.Response {
	if len(service.IDs) == 0 {
		return serializer.Response{
			Status: 400,
			Msg:    "请传入要删除的ID列表",
		}
	}

	if err := model.DB.Delete(&model.Classmate{}, service.IDs).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "删除失败",
		}
	}

	return serializer.Response{
		Status: 200,
		Msg:    "批量删除成功",
	}
}
