package service

import (
	"sort"
	"strconv"
	"strings"
	"todo_list/cache"
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

	// 更新缓存
	cache.UpdateClassmate(classmate)

	return serializer.Response{Status: 200, Msg: "创建成功"}
}

func (service *ListClassmateService) List(uid uint) serializer.Response {
	var results []model.Classmate

	cache.ClassmateCacheLock.RLock()
	for _, c := range cache.ClassmateCache {
		if c.Uid != uid {
			continue
		}
		if service.Name != "" && !strings.Contains(strings.ToLower(c.Name), strings.ToLower(service.Name)) {
			continue
		}
		results = append(results, c)
	}
	cache.ClassmateCacheLock.RUnlock()

	total := len(results)

	// 按创建时间倒序
	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.After(results[j].CreatedAt)
	})

	// 分页
	if service.Page <= 0 {
		service.Page = 1
	}
	if service.PageSize <= 0 {
		service.PageSize = 10
	}
	start := (service.Page - 1) * service.PageSize
	end := start + service.PageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	paged := results[start:end]

	return serializer.Response{
		Status: 200,
		Msg:    "获取同学列表成功",
		Data: map[string]interface{}{
			"list":      paged,
			"total":     total,
			"page":      service.Page,
			"page_size": service.PageSize,
		},
	}
}

func (service *UpdateClassmateService) Update(idStr string) serializer.Response {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return serializer.Response{Status: 400, Msg: "ID无效"}
	}

	var classmate model.Classmate
	if err := model.DB.First(&classmate, id).Error; err != nil {
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

	// 更新缓存
	cache.UpdateClassmate(classmate)

	return serializer.Response{Status: 200, Msg: "更新成功"}
}

// 批量删除
func (service *BatchDeleteClassmateService) BatchDelete() serializer.Response {
	if len(service.IDs) == 0 {
		return serializer.Response{Status: 400, Msg: "请传入要删除的ID列表"}
	}

	if err := model.DB.Delete(&model.Classmate{}, service.IDs).Error; err != nil {
		return serializer.Response{Status: 500, Msg: "删除失败"}
	}

	// 删除缓存
	for _, id := range service.IDs {
		cache.DeleteClassmate(id)
	}

	return serializer.Response{Status: 200, Msg: "批量删除成功"}
}

func (service *ListClassmateService) GetUnsetMapClassmates(id uint) serializer.Response {
	var classmates []model.Classmate

	// 先从缓存中筛选
	cache.ClassmateCacheLock.RLock()
	if len(cache.ClassmateCache) > 0 {
		for _, classmate := range cache.ClassmateCache {
			if classmate.IsSetMap == 0 { // 筛选未设置地图点位的同学
				classmates = append(classmates, classmate)
			}
		}
		cache.ClassmateCacheLock.RUnlock()
	} else {
		// 缓存为空，从数据库查
		cache.ClassmateCacheLock.RUnlock()
		if err := model.DB.Where("is_set_map = ?", 0).Find(&classmates).Error; err != nil {
			return serializer.Response{
				Status: 500,
				Msg:    "获取未设置地图同学信息失败: " + err.Error(),
			}
		}

		// 同时把查到的所有（含 is_set_map=0 和 ≠0）加载进缓存
		var allClassmates []model.Classmate
		if err := model.DB.Find(&allClassmates).Error; err == nil {
			cache.ClassmateCacheLock.Lock()
			for _, c := range allClassmates {
				cache.ClassmateCache[c.ID] = c
			}
			cache.ClassmateCacheLock.Unlock()
		}
	}

	simpleClassmates := serializer.BuildClassmatesMapView(classmates)

	return serializer.Response{
		Status: 200,
		Data:   simpleClassmates,
		Msg:    "获取成功",
	}
}
