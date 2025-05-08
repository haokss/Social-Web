package model

import "github.com/jinzhu/gorm"

type Friend struct {
	Uid      uint   `json:"uid"`      // 所属用户 ID
	Name     string `json:"name"`     // 姓名
	Phone    string `json:"phone"`    // 联系方式
	Birthday string `json:"birthday"` // 生日
	Hobby    string `json:"hobby"`    // 兴趣爱好
	gorm.Model
}
