package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// 计划任务
type Image struct {
	User        User   `gorm:"ForeignKey:Uid"`
	Uid         uint   `gorm:"not null"`
	ImageName   string `gorm:"index; not null"`
	Url         string `gorm:"not null"`
	Description string `gorm:""`
	IsChecked   int    `gorm:"default:0"`
	CreateTime  time.Time
	gorm.Model
}
