package model

type Point struct {
	ID             uint    `json:"id" gorm:"primaryKey"`
	UserID         uint    `json:"user_id" gorm:"not null"`
	Name           string  `json:"name" gorm:"not null"`
	Type           string  `json:"type" gorm:"not null"`
	Latitude       float64 `json:"lat"`
	Longitude      float64 `json:"lng"`
	SelectedPeople string  `json:"selected_people"`
	CreatedAt      int64   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      int64   `json:"updated_at" gorm:"autoUpdateTime"`
}
