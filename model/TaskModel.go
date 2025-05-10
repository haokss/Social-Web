package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// 计划任务 计划任务的最小单位为天，可以没有结束时间，设置结束时间就认为可以逾期
type Task struct {
	User      User   `gorm:"ForeignKey:Uid"`
	Uid       uint   `gorm:"not null"`
	Title     string `gorm:"index; not null"`
	Status    int    `gorm:"default:'0'"` // 0未完成，1 已完成
	Content   string `gorm:"type:longtext"`
	Priority  int    `grom:"default:'5'"` // 任务优先级：0~5
	IsNotify  int    `gorm:"default:'0'"` // 0未完成，1 已完成
	NotifyWay int    `gorm:"default:0"`   // 新增：提醒方式（0站内信/1邮件）
	IsChecked int    `gorm:"default:0"`
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
	IsChecked int       `gorm:"default:0"`
	StartTime time.Time `gorm:"type:datetime"`
	EndTime   time.Time `gorm:"type:datetime"`
	gorm.Model
}
