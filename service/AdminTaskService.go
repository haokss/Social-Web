package service

import (
	"todo_list/cache"
	"todo_list/model"
	"todo_list/serializer"

	"github.com/jinzhu/gorm"
)

type AdminTaskListService struct{}

type AuditItem struct {
	TargetUserID uint   `json:"target_user_id"`
	TaskID       []uint `json:"task_id"`
}

type AdminBatchAuditService struct {
	IsChecked  uint        `json:"is_checked" binding:"required"`
	AuditItems []AuditItem `json:"audit_items" binding:"required"`
}

// 获取所有用户的任务
func (service *AdminTaskListService) GetAllUserTasks() serializer.Response {
	tasks := cache.GetAllTasks()

	// 用 Preload 预加载关联的 User
	if err := model.DB.Preload("User").Find(&tasks).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "获取任务失败",
		}
	}

	// 构造返回的任务列表数据
	var taskList []serializer.AdminTaskSerializer
	for _, task := range tasks {
		taskList = append(taskList, serializer.AdminBuildTask(task))
	}

	return serializer.Response{
		Status: 200,
		Msg:    "获取所有用户任务成功",
		Data: map[string]interface{}{
			"tasks": taskList,
		},
	}
}

func (service *AdminBatchAuditService) BatchAudit() serializer.Response {
	db := model.DB

	// 用事务确保一致性
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, item := range service.AuditItems {
			if len(item.TaskID) == 0 {
				continue
			}

			var tasks []model.Task
			// 查询要更新的任务
			if err := tx.Debug().Where("uid = ? AND id IN (?)", item.TargetUserID, item.TaskID).Find(&tasks).Error; err != nil {
				return err
			}

			// 更新数据库
			if err := tx.Debug().Model(&model.Task{}).
				Where("uid = ? AND id IN (?)", item.TargetUserID, item.TaskID).
				Update("is_checked", service.IsChecked).Error; err != nil {
				return err
			}

			// 更新缓存
			for _, task := range tasks {
				// task.IsChecked = service.IsChecked
				cache.SetTask(task)
			}
		}
		return nil
	})

	// TODO：审核流程通知到用户

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
