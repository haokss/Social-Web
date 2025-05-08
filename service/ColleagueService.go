package service

import (
	"todo_list/cache"
	"todo_list/model"
	"todo_list/package/utils"
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
	// 写入缓存
	cache.UpdateColleagueCache(colleague)
	return serializer.Response{
		Status: 200,
		Msg:    "创建成功",
	}
}

func (service *ListColleagueService) List(uid uint) serializer.Response {
	var result []model.Colleague

	// 从缓存读取所有同事
	cache.ColleagueCacheLock.RLock()
	for _, c := range cache.ColleagueCache {
		if c.Uid != uid {
			continue
		}
		if service.Name != "" && !utils.ContainsIgnoreCase(c.Name, service.Name) {
			continue
		}
		if service.Company != "" && !utils.ContainsIgnoreCase(c.Company, service.Company) {
			continue
		}
		result = append(result, c)
	}
	cache.ColleagueCacheLock.RUnlock()

	total := len(result)

	// 分页处理
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
	pagedResult := result[start:end]

	return serializer.Response{
		Status: 200,
		Msg:    "获取同事列表成功（缓存）",
		Data: map[string]interface{}{
			"list":      pagedResult,
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

	// 更新缓存
	cache.UpdateColleagueCache(colleague)

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

	// 删除缓存
	cache.DeleteColleagueCache(service.IDs)

	return serializer.Response{
		Status: 200,
		Msg:    "批量删除成功",
	}
}

func (service *ListColleagueService) GetUnsetMapColleagues(id uint) serializer.Response {
	var colleagues []model.Colleague

	// 先从缓存中筛选
	cache.ColleagueCacheLock.RLock()
	if len(cache.ColleagueCache) > 0 {
		for _, colleague := range cache.ColleagueCache {
			if colleague.IsSetMap == 0 { // 筛选未设置点位的同事
				colleagues = append(colleagues, colleague)
			}
		}
		cache.ColleagueCacheLock.RUnlock()
	} else {
		// 缓存为空，从数据库查
		cache.ColleagueCacheLock.RUnlock()
		if err := model.DB.Where("is_set_map = ?", 0).Find(&colleagues).Error; err != nil {
			return serializer.Response{
				Status: 500,
				Msg:    "获取未设置地图同事信息失败: " + err.Error(),
			}
		}

		// 同时把查到的所有（含 is_set_map=0 和 ≠0）加载进缓存
		var allColleagues []model.Colleague
		if err := model.DB.Find(&allColleagues).Error; err == nil {
			cache.ColleagueCacheLock.Lock()
			for _, c := range allColleagues {
				cache.ColleagueCache[c.ID] = c
			}
			cache.ColleagueCacheLock.Unlock()
		}
	}

	simpleColleagues := serializer.BuildColleaguesMapView(colleagues)

	return serializer.Response{
		Status: 200,
		Data:   simpleColleagues,
		Msg:    "获取成功",
	}
}
