package model

// 家庭成员关系
type RelativeInfo struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"type:varchar(100)" json:"name"`
	Relation        string         `gorm:"type:varchar(50)" json:"relation"`
	Gender          string         `gorm:"type:varchar(10)" json:"gender"`
	Address         string         `gorm:"type:text" json:"address"`
	Contact         string         `gorm:"type:varchar(100)" json:"contact"`
	WeChat          string         `gorm:"type:varchar(100)" json:"wechat"`
	HasDebtRelation bool           `gorm:"default:false" json:"hasDebtRelation"`
	DebtType        string         `gorm:"type:varchar(50)" json:"debtType"`
	DebtProof       string         `gorm:"type:text" json:"debtProof"`
	Note            string         `gorm:"type:text" json:"note"`
	Avatar          string         `gorm:"type:text" json:"avatar"`
	ParentID        *uint          `json:"parentId"` // Nullable for root
	Children        []RelativeInfo `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}
