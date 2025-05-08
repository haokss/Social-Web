package cache

import (
	"sync"
	"todo_list/model"
)

// 同学缓存，key: classmate ID
var (
	ClassmateCache     = make(map[uint]model.Classmate)
	ClassmateCacheLock sync.RWMutex
)

// 初始化缓存：加载所有同学
func InitClassmateCache() error {
	var classmates []model.Classmate
	if err := model.DB.Find(&classmates).Error; err != nil {
		return err
	}
	ClassmateCacheLock.Lock()
	defer ClassmateCacheLock.Unlock()
	for _, c := range classmates {
		ClassmateCache[c.ID] = c
	}
	return nil
}

// 获取单个同学
func GetClassmate(id uint) (model.Classmate, bool) {
	ClassmateCacheLock.RLock()
	defer ClassmateCacheLock.RUnlock()
	classmate, ok := ClassmateCache[id]
	return classmate, ok
}

// 更新或新增同学到缓存
func UpdateClassmate(classmate model.Classmate) {
	ClassmateCacheLock.Lock()
	defer ClassmateCacheLock.Unlock()
	ClassmateCache[classmate.ID] = classmate
}

// 删除缓存中的同学
func DeleteClassmate(id uint) {
	ClassmateCacheLock.Lock()
	defer ClassmateCacheLock.Unlock()
	delete(ClassmateCache, id)
}
