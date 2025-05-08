package model

import (
	"github.com/jinzhu/gorm"
)

type Classmate struct {
	Uid       uint   `json:"uid"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	School    string `json:"school"`
	ClassName string `json:"className"`
	Stage     string `json:"stage"`
	Birthday  string `json:"birthday"`
	Hobby     string `json:"hobby"`
	gorm.Model
}
