package model

import (
	"github.com/jinzhu/gorm"
)

type Bill struct {
	TransactionTime string  `gorm:"size:100"`
	TransactionType string  `gorm:"size:50;not null"`
	Counterparty    string  `gorm:"size:100"` // 允许为空
	Product         string  `gorm:"size:100"`
	IncomeExpense   string  `gorm:"size:10;not null"`            // 收入/支出
	Amount          float64 `gorm:"type:decimal(12,2);not null"` // 精确小数
	PaymentMethod   string  `gorm:"size:30;not null"`
	Status          string  `gorm:"size:20;not null"`
	TransactionID   string  `gorm:"size:64;uniqueIndex;not null"` // 防重复
	MerchantID      string  `gorm:"size:64"`
	Remarks         string  `gorm:"type:text"`
	gorm.Model
}
