package cache

import (
	"sync"
	"todo_list/model"
)

var (
	ColleagueCache     = make(map[uint]model.Colleague) // key: colleague ID
	ColleagueCacheLock sync.RWMutex
)

// InitColleagueCache 在服务启动时调用，加载所有同事到内存
func InitColleagueCache() error {
	var colleagues []model.Colleague
	if err := model.DB.Find(&colleagues).Error; err != nil {
		return err
	}
	ColleagueCacheLock.Lock()
	defer ColleagueCacheLock.Unlock()
	for _, c := range colleagues {
		ColleagueCache[c.ID] = c
	}
	return nil
}

// GetColleague 从缓存中获取同事
func GetColleague(id uint) (model.Colleague, bool) {
	ColleagueCacheLock.RLock()
	defer ColleagueCacheLock.RUnlock()
	c, ok := ColleagueCache[id]
	return c, ok
}

// UpdateColleagueCache 更新/添加缓存
func UpdateColleagueCache(c model.Colleague) {
	ColleagueCacheLock.Lock()
	defer ColleagueCacheLock.Unlock()
	ColleagueCache[c.ID] = c
}

// DeleteColleagueCache 批量删除缓存
func DeleteColleagueCache(ids []uint) {
	ColleagueCacheLock.Lock()
	defer ColleagueCacheLock.Unlock()
	for _, id := range ids {
		delete(ColleagueCache, id)
	}
}
