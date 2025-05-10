package service

import (
	"todo_list/cache"
	"todo_list/model"
	"todo_list/serializer"

	"github.com/jinzhu/gorm"
)

type AdminTimingTaskListService struct{}

type AuditTimingItem struct {
	TargetUserID uint   `json:"target_user_id"`
	TaskID       []uint `json:"task_id"`
}

type AdminBatchAuditTimingTaskService struct {
	AuditItems []AuditTimingItem `json:"audit_items"`
	IsChecked  int               `json:"is_checked"` // 审核状态
}

func (service *AdminTimingTaskListService) GetAllTimingTasks() serializer.Response {
	var timingTasks []model.TimingTask

	// 用 Preload 预加载关联的 User
	if err := model.DB.Preload("User").Find(&timingTasks).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "获取定时任务失败",
		}
	}

	// 构造返回的定时任务列表数据
	var timingTaskList []serializer.AdminTimingTaskSerializer
	for _, task := range timingTasks {
		timingTaskList = append(timingTaskList, serializer.AdminBuildTimingTask(task))
	}

	return serializer.Response{
		Status: 200,
		Msg:    "获取所有用户定时任务成功",
		Data: map[string]interface{}{
			"timing_tasks": timingTaskList,
		},
	}
}

func (service *AdminBatchAuditTimingTaskService) BatchAuditTimingTask() serializer.Response {
	db := model.DB

	err := db.Transaction(func(tx *gorm.DB) error {
		for _, item := range service.AuditItems {
			if len(item.TaskID) == 0 {
				continue
			}

			var tasks []model.TimingTask
			// 查询需要更新的定时任务
			if err := tx.Where("uid = ? AND id IN (?)", item.TargetUserID, item.TaskID).Find(&tasks).Error; err != nil {
				return err
			}

			// 更新数据库 is_checked 字段
			if err := tx.Model(&model.TimingTask{}).
				Where("uid = ? AND id IN (?)", item.TargetUserID, item.TaskID).
				Update("is_checked", service.IsChecked).Error; err != nil {
				return err
			}

			// 更新缓存（如果有的话）
			for _, task := range tasks {
				task.IsChecked = service.IsChecked
				cache.SetTimingTask(task) // 伪代码，根据你的实际缓存实现调整
			}
		}
		return nil
	})

	// TODO: 可在这里添加审核结果通知

	if err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "批量审核失败",
		}
	}

	return serializer.Response{
		Status: 200,
		Msg:    "批量审核成功",
	}
}
