package service

import (
	"todo_list/model"
	"todo_list/serializer"
)

type CreateColleagueService struct {
	Name     string `json:"name" form:"name"`
	Company  string `json:"company" form:"company"`
	Position string `json:"position" form:"position"`
	Phone    string `json:"phone" form:"phone"`
	Email    string `json:"email" form:"email"`
}

type ListColleagueService struct {
	Name     string `form:"name"`
	Company  string `form:"company"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type UpdateColleagueService struct {
	Name     string `json:"name" form:"name"`
	Company  string `json:"company" form:"company"`
	Position string `json:"position" form:"position"`
	Phone    string `json:"phone" form:"phone"`
	Email    string `json:"email" form:"email"`
}

type BatchDeleteColleagueService struct {
	IDs []uint `json:"ids" form:"ids"`
}

func (service *CreateColleagueService) Create(uid uint) serializer.Response {
	colleague := model.Colleague{
		Uid:      uid,
		Name:     service.Name,
		Company:  service.Company,
		Position: service.Position,
		Phone:    service.Phone,
		Email:    service.Email,
	}
	if err := model.DB.Create(&colleague).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "创建失败",
		}
	}
	return serializer.Response{
		Status: 200,
		Msg:    "创建成功",
	}
}

func (service *ListColleagueService) List(uid uint) serializer.Response {
	var colleagues []model.Colleague
	var total int64

	if service.Page <= 0 {
		service.Page = 1
	}
	if service.PageSize <= 0 {
		service.PageSize = 10
	}

	tx := model.DB.Model(&model.Colleague{}).Where("uid = ?", uid)
	if service.Name != "" {
		tx = tx.Where("name LIKE ?", "%"+service.Name+"%")
	}
	if service.Company != "" {
		tx = tx.Where("company LIKE ?", "%"+service.Company+"%")
	}

	tx.Count(&total)

	err := tx.Order("created_at desc").
		Offset((service.Page - 1) * service.PageSize).
		Limit(service.PageSize).
		Find(&colleagues).Error

	if err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "获取同事列表失败",
		}
	}

	return serializer.Response{
		Status: 200,
		Msg:    "获取同事列表成功",
		Data: map[string]interface{}{
			"list":      colleagues,
			"total":     total,
			"page":      service.Page,
			"page_size": service.PageSize,
		},
	}
}

func (service *UpdateColleagueService) Update(idStr string) serializer.Response {
	var colleague model.Colleague
	if err := model.DB.First(&colleague, idStr).Error; err != nil {
		return serializer.Response{
			Status: 404,
			Msg:    "同事未找到",
		}
	}

	colleague.Name = service.Name
	colleague.Company = service.Company
	colleague.Position = service.Position
	colleague.Phone = service.Phone
	colleague.Email = service.Email

	if err := model.DB.Save(&colleague).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "更新失败",
		}
	}
	return serializer.Response{
		Status: 200,
		Msg:    "更新成功",
	}
}

func (service *BatchDeleteColleagueService) BatchDelete() serializer.Response {
	if len(service.IDs) == 0 {
		return serializer.Response{
			Status: 400,
			Msg:    "请传入要删除的ID列表",
		}
	}

	if err := model.DB.Delete(&model.Colleague{}, service.IDs).Error; err != nil {
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
