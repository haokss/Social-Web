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
	IsSetMap  uint   `json:"is_set_map"` // 是否设置地图点位
	MapUid    uint   `json:"map_uid"`    // 地图点位ID
	gorm.Model
}
