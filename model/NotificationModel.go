package model

import (
	"time"
)

type Notification struct {
	ID         uint   `gorm:"primarykey"`
	UserID     uint   `gorm:"index;not null"`
	Event      string `gorm:"size:100;not null"`
	Data       string `gorm:"type:text;not null"`
	ScheduleID string `gorm:"size:255"`
	IsRead     bool   `gorm:"default:false"`
	CreatedAt  time.Time
}
