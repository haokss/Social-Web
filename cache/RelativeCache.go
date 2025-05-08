package cache

import (
	"sync"
	"todo_list/model"
)

// 定义亲属缓存
var (
	RelativeCache     = make(map[uint]model.RelativeInfo) // key: relative ID
	RelativeCacheLock sync.RWMutex
)

// InitRelativeCache 在服务启动时调用，加载所有亲属信息到内存
func InitRelativeCache() error {
	var relatives []model.RelativeInfo
	if err := model.DB.Find(&relatives).Error; err != nil {
		return err
	}
	RelativeCacheLock.Lock()
	defer RelativeCacheLock.Unlock()
	for _, relative := range relatives {
		RelativeCache[relative.ID] = relative
	}
	return nil
}

// GetRelative 从缓存中获取亲属
func GetRelative(id uint) (model.RelativeInfo, bool) {
	RelativeCacheLock.RLock()
	defer RelativeCacheLock.RUnlock()
	relative, ok := RelativeCache[id]
	return relative, ok
}

// UpdateRelative 更新缓存中的亲属信息
func UpdateRelative(relative model.RelativeInfo) {
	RelativeCacheLock.Lock()
	defer RelativeCacheLock.Unlock()
	RelativeCache[relative.ID] = relative
}

// DeleteRelative 从缓存中删除亲属信息
func DeleteRelative(id uint) {
	RelativeCacheLock.Lock()
	defer RelativeCacheLock.Unlock()
	delete(RelativeCache, id)
}
