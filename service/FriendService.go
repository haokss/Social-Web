package service

import (
	"sort"
	"strconv"
	"strings"
	"todo_list/cache"
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

	// 更新缓存
	cache.UpdateFriend(friend)

	return serializer.Response{Status: 200, Msg: "创建成功"}
}

// 查询朋友列表
func (service *ListFriendService) List(uid uint) serializer.Response {
	var results []model.Friend

	cache.FriendCacheLock.RLock()
	for _, f := range cache.FriendCache {
		if f.Uid != uid {
			continue
		}
		if service.Name != "" && !strings.Contains(strings.ToLower(f.Name), strings.ToLower(service.Name)) {
			continue
		}
		results = append(results, f)
	}
	cache.FriendCacheLock.RUnlock()

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
		Msg:    "获取朋友列表成功",
		Data: map[string]interface{}{
			"list":      paged,
			"total":     total,
			"page":      service.Page,
			"page_size": service.PageSize,
		},
	}
}

// 更新朋友
func (service *UpdateFriendService) Update(idStr string) serializer.Response {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return serializer.Response{Status: 400, Msg: "ID无效"}
	}

	var friend model.Friend
	if err := model.DB.First(&friend, id).Error; err != nil {
		return serializer.Response{Status: 404, Msg: "朋友未找到"}
	}

	friend.Name = service.Name
	friend.Phone = service.Phone
	friend.Birthday = service.Birthday
	friend.Hobby = service.Hobby

	if err := model.DB.Save(&friend).Error; err != nil {
		return serializer.Response{Status: 500, Msg: "更新失败"}
	}

	// 更新缓存
	cache.UpdateFriend(friend)

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

	// 删除缓存
	for _, id := range service.IDs {
		cache.DeleteFriend(id)
	}

	return serializer.Response{Status: 200, Msg: "批量删除成功"}
}

func (service *ListFriendService) GetUnsetMapFriends(id uint) serializer.Response {
	var friends []model.Friend

	// 先从缓存中筛选
	cache.FriendCacheLock.RLock()
	if len(cache.FriendCache) > 0 {
		for _, friend := range cache.FriendCache {
			if friend.IsSetMap == 0 { // 筛选未设置地图点位的朋友
				friends = append(friends, friend)
			}
		}
		cache.FriendCacheLock.RUnlock()
	} else {
		// 缓存为空，从数据库查
		cache.FriendCacheLock.RUnlock()
		if err := model.DB.Where("is_set_map = ?", 0).Find(&friends).Error; err != nil {
			return serializer.Response{
				Status: 500,
				Msg:    "获取未设置地图朋友信息失败: " + err.Error(),
			}
		}

		// 同时把查到的所有（含 is_set_map=0 和 ≠0）加载进缓存
		var allFriends []model.Friend
		if err := model.DB.Find(&allFriends).Error; err == nil {
			cache.FriendCacheLock.Lock()
			for _, f := range allFriends {
				cache.FriendCache[f.ID] = f
			}
			cache.FriendCacheLock.Unlock()
		}
	}

	// 构建返回值
	simpleFriends := serializer.BuildFriendsMapView(friends)
	return serializer.Response{
		Status: 200,
		Data:   simpleFriends,
		Msg:    "获取成功",
	}
}
