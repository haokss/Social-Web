package model

import "github.com/jinzhu/gorm"

type Colleague struct {
	Uid      uint   `json:"uid"`      // 所属用户ID
	Name     string `json:"name"`     // 姓名
	Company  string `json:"company"`  // 工作单位
	Position string `json:"position"` // 职位
	Phone    string `json:"phone"`    // 手机号
	Email    string `json:"email"`    // 邮箱
	gorm.Model
}
