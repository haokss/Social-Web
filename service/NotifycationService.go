package service

import (
	"todo_list/model"
	"todo_list/serializer"
)

type NotificationListService struct{}

func (service *NotificationListService) List(uid uint) serializer.Response {
	var notifications []model.Notification

	if err := model.DB.
		Where("user_id = ?", uid).
		Order("created_at desc").
		Limit(100).
		Find(&notifications).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "获取通知失败",
		}
	}

	return serializer.Response{
		Status: 200,
		Data:   serializer.BuildNotificationList(notifications),
		Msg:    "获取成功",
	}
}

func (service *NotificationListService) MarkNotificationRead(uid uint, id string) serializer.Response {
	var notification model.Notification

	// 查找指定消息并标记为已读
	if err := model.DB.Model(&notification).
		Where("id = ? AND user_id = ?", id, uid).
		Update("is_read", true).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "标记失败",
		}
	}

	return serializer.Response{
		Status: 200,
		Msg:    "标记为已读成功",
	}
}

func (service *NotificationListService) MarkAllNotificationsRead(uid uint) serializer.Response {
	// 一键标记所有未读消息为已读
	if err := model.DB.Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", uid, false).
		Update("is_read", true).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "一键已读失败",
		}
	}

	return serializer.Response{
		Status: 200,
		Msg:    "全部标记为已读成功",
	}
}
