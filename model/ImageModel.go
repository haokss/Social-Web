package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// 计划任务
type Image struct {
	Uid         uint   `gorm:"not null"`
	ImageName   string `gorm:"index; not null"`
	Url         string `gorm:"not null"`
	Description string `gorm:""`
	CreateTime  time.Time
	gorm.Model
}
