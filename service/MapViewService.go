package service

import (
	"fmt"
	"strings"
	"todo_list/model"
	"todo_list/serializer"
)

// CreatePointService 用于接收和处理点位创建请求的数据
type CreatePointService struct {
	Name           string `json:"name" binding:"required"`
	Type           string `json:"type" binding:"required"`
	LatLng         LatLng `json:"latlng" binding:"required"`
	SelectedPeople []uint `json:"selectedPeople" binding:"required"`
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type ListPointService struct{}

// Create 创建点位并保存到数据库
func (service *CreatePointService) Create(userId uint) serializer.Response {
	var idStrs []string
	for _, id := range service.SelectedPeople {
		idStrs = append(idStrs, fmt.Sprintf("%d", id))
	}
	selectedPeopleStr := strings.Join(idStrs, ",")

	point := model.Point{
		UserID:         userId,
		Name:           service.Name,
		Type:           service.Type,
		Latitude:       service.LatLng.Lat,
		Longitude:      service.LatLng.Lng,
		SelectedPeople: selectedPeopleStr,
	}

	if err := model.DB.Create(&point).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "创建点位失败: " + err.Error(),
		}
	}

	return serializer.Response{
		Status: 200,
		Msg:    "创建点位成功",
		Data:   point,
	}
}

func (service *ListPointService) List(userID uint) serializer.Response {
	var points []model.Point
	if err := model.DB.Where("user_id = ?", userID).Find(&points).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "查询点位失败: " + err.Error(),
		}
	}

	return serializer.Response{
		Status: 200,
		Msg:    "查询点位成功",
		Data:   points,
	}
}

// Update 更新点位
func (service *ListPointService) Update(point model.Point, userID uint) serializer.Response {
	var existingPoint model.Point

	// 查找指定 ID 的点位
	if err := model.DB.Where("id = ? AND user_id = ?", point.ID, userID).First(&existingPoint).Error; err != nil {
		return serializer.Response{
			Status: 404,
			Msg:    "点位未找到",
		}
	}

	// 更新点位数据
	existingPoint.Name = point.Name
	existingPoint.Type = point.Type
	existingPoint.Latitude = point.Latitude
	existingPoint.Longitude = point.Longitude
	existingPoint.SelectedPeople = point.SelectedPeople

	// 保存更新后的点位
	if err := model.DB.Save(&existingPoint).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "更新点位失败: " + err.Error(),
		}
	}

	return serializer.Response{
		Status: 200,
		Msg:    "点位更新成功",
		Data:   existingPoint,
	}
}

// Delete 删除点位
func (service *ListPointService) Delete(pointID string, userID uint) serializer.Response {
	var point model.Point

	// 查找指定 ID 的点位
	if err := model.DB.Where("id = ? AND user_id = ?", pointID, userID).First(&point).Error; err != nil {
		return serializer.Response{
			Status: 404,
			Msg:    "点位未找到",
		}
	}

	// 删除点位
	if err := model.DB.Delete(&point).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "删除点位失败: " + err.Error(),
		}
	}

	return serializer.Response{
		Status: 200,
		Msg:    "点位删除成功",
	}
}
