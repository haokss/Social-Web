package cache

import (
	"sync"
	"todo_list/model"
)

// 朋友缓存，key: friend ID
var (
	FriendCache     = make(map[uint]model.Friend)
	FriendCacheLock sync.RWMutex
)

// 初始化缓存：加载所有朋友
func InitFriendCache() error {
	var friends []model.Friend
	if err := model.DB.Find(&friends).Error; err != nil {
		return err
	}
	FriendCacheLock.Lock()
	defer FriendCacheLock.Unlock()
	for _, f := range friends {
		FriendCache[f.ID] = f
	}
	return nil
}

// 获取单个朋友
func GetFriend(id uint) (model.Friend, bool) {
	FriendCacheLock.RLock()
	defer FriendCacheLock.RUnlock()
	friend, ok := FriendCache[id]
	return friend, ok
}

// 更新或新增朋友到缓存
func UpdateFriend(friend model.Friend) {
	FriendCacheLock.Lock()
	defer FriendCacheLock.Unlock()
	FriendCache[friend.ID] = friend
}

// 删除缓存中的朋友
func DeleteFriend(id uint) {
	FriendCacheLock.Lock()
	defer FriendCacheLock.Unlock()
	delete(FriendCache, id)
}
