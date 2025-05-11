package service

import (
	"fmt"
	"strings"
	"todo_list/cache"
	"todo_list/model"
	sse "todo_list/package/SSE"
	"todo_list/serializer"

	"github.com/gin-gonic/gin"
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

func (service *AdminBatchAuditService) BatchAudit(c *gin.Context) serializer.Response {
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

	// TODO：更新审核状态到缓存

	// 审核流程状态通知到用户
	broker := c.MustGet("sseBroker").(*sse.Broker)
	for _, item := range service.AuditItems {
		if len(item.TaskID) == 0 {
			continue
		}

		// 重新查询任务标题（或缓存 tasks 数组供此处使用）
		var tasks []model.Task
		if err := model.DB.Where("uid = ? AND id IN (?)", item.TargetUserID, item.TaskID).Find(&tasks).Error; err != nil {
			// 忽略通知失败
			continue
		}

		// 提取活动名称
		var titles []string
		for _, task := range tasks {
			titles = append(titles, task.Title)
		}

		statusStr := "审核通过"
		if service.IsChecked == 2 {
			statusStr = "审核未通过"
		}

		content := fmt.Sprintf("您提交的活动【%s】已被管理员%s", strings.Join(titles, "、"), statusStr)

		data := map[string]interface{}{
			"title":    "审核提醒",
			"content":  content,
			"task_ids": item.TaskID,
		}

		broker.Notify(sse.Message{
			Event:     "instant_notification",
			Data:      data,
			TargetIDs: []uint{item.TargetUserID},
		})
	}

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
