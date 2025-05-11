package serializer

import "todo_list/model"

type Notification struct {
	ID        uint        `json:"id"`
	Event     string      `json:"event"`
	Data      interface{} `json:"data"` // 若你存储为 JSON，可以用 datatypes.JSON 原样传出
	IsRead    bool        `json:"is_read"`
	CreatedAt int64       `json:"created_at"`
}

func BuildNotification(item model.Notification) Notification {
	return Notification{
		ID:        item.ID,
		Event:     item.Event,
		Data:      item.Data,
		IsRead:    item.IsRead,
		CreatedAt: item.CreatedAt.Unix(),
	}
}

func BuildNotificationList(items []model.Notification) []Notification {
	var list []Notification
	for _, item := range items {
		list = append(list, BuildNotification(item))
	}
	return list
}
