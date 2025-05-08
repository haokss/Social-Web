package cache

import (
	"sync"
	"todo_list/model"
)

var (
	TaskCache     = make(map[uint]model.Task) // key: task ID
	TaskCacheLock sync.RWMutex
)

// InitTaskCache 在服务启动时调用，加载所有任务到内存
func InitTaskCache() error {
	var tasks []model.Task
	if err := model.DB.Find(&tasks).Error; err != nil {
		return err
	}
	TaskCacheLock.Lock()
	defer TaskCacheLock.Unlock()
	for _, task := range tasks {
		TaskCache[task.ID] = task
	}
	// log.Info("[TaskCache] Loaded %d tasks into memory", len(TaskCache))
	return nil
}

// GetTask 从缓存中获取任务
func GetTask(tid uint) (model.Task, bool) {
	TaskCacheLock.RLock()
	defer TaskCacheLock.RUnlock()
	task, ok := TaskCache[tid]
	return task, ok
}

// SetTask 更新或添加任务到缓存
func SetTask(task model.Task) {
	TaskCacheLock.Lock()
	defer TaskCacheLock.Unlock()
	TaskCache[task.ID] = task
}

// DeleteTask 从缓存中删除任务
func DeleteTask(tid uint) {
	TaskCacheLock.Lock()
	defer TaskCacheLock.Unlock()
	delete(TaskCache, tid)
}

// GetUserTasks 获取指定用户的所有任务
func GetUserTasks(uid uint) []model.Task {
	TaskCacheLock.RLock()
	defer TaskCacheLock.RUnlock()
	var result []model.Task
	for _, task := range TaskCache {
		if task.Uid == uid {
			result = append(result, task)
		}
	}
	return result
}

// GetUserTaskStats 获取指定用户的任务统计信息
func GetUserTaskStats(uid uint) (total int, done int) {
	TaskCacheLock.RLock()
	defer TaskCacheLock.RUnlock()
	for _, task := range TaskCache {
		if task.Uid == uid {
			total++
			if task.Status == 1 {
				done++
			}
		}
	}
	return
}
