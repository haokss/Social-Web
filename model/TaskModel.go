package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// 计划任务
type Task struct {
	User      User   `gorm:"ForeignKey:Uid"`
	Uid       uint   `gorm:"not null"`
	Title     string `gorm:"index; not null"`
	Status    int    `gorm:"default:'0'"` // 0未完成，1 已完成
	Content   string `gorm:"type:longtext"`
	StartTime int64
	EndTime   int64
	gorm.Model
}

// 定时任务
type TimingTask struct {
	User      User      `gorm:"ForeignKey:Uid"`
	Uid       uint      `gorm:"not null"`
	Title     string    `gorm:"index; not null"`
	Status    int       `gorm:"default:'0'"` // 0未完成，1 已完成
	Content   string    `gorm:"type:longtext"`
	Priority  int       `grom:"default:'5'"` // 任务优先级：0~5
	Type      int       `gorm:"default:'0'"` // 任务类型
	NotifyWay int       `gorm:"default:0"`   // 新增：提醒方式（0站内信/1邮件/2短信）
	EarlyTime int       `gorm:"default:0"`   // 提前的时间，秒
	StartTime time.Time `gorm:"type:datetime"`
	EndTime   time.Time `gorm:"type:datetime"`
	gorm.Model
}
