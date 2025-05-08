package cache

import (
	"sync"
	"todo_list/model"
)

var (
	TimingTaskCache     = make(map[uint]model.TimingTask) // key: task ID
	TimingTaskCacheLock sync.RWMutex
)

// 初始化缓存：从数据库加载任务
func InitTimingTaskCache() error {
	var allTasks []model.TimingTask
	if err := model.DB.Find(&allTasks).Error; err != nil {
		return err
	}
	TimingTaskCacheLock.Lock()
	defer TimingTaskCacheLock.Unlock()
	for _, task := range allTasks {
		TimingTaskCache[task.ID] = task
	}
	// log.Printf("TimingTask cache initialized with %d tasks", len(allTasks))
	return nil
}

// 添加或更新任务到缓存
func SetTimingTask(task model.TimingTask) {
	TimingTaskCacheLock.Lock()
	defer TimingTaskCacheLock.Unlock()
	TimingTaskCache[task.ID] = task
}

// 根据 ID 获取任务（存在则返回 true）
func GetTimingTask(id uint) (model.TimingTask, bool) {
	TimingTaskCacheLock.RLock()
	defer TimingTaskCacheLock.RUnlock()
	task, exists := TimingTaskCache[id]
	return task, exists
}

// 删除缓存中的任务
func DeleteTimingTask(id uint) {
	TimingTaskCacheLock.Lock()
	defer TimingTaskCacheLock.Unlock()
	delete(TimingTaskCache, id)
}

// 清空整个缓存
func ClearTimingTaskCache() {
	TimingTaskCacheLock.Lock()
	defer TimingTaskCacheLock.Unlock()
	TimingTaskCache = make(map[uint]model.TimingTask)
}
