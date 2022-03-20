package models

import "time"

// BaseModel is a base model for all models
type BaseModel struct {
	Deleted   bool      `json:"is_deleted" gorm:"default:false;not null"`
	CreatedAt time.Time `json:"-" gorm:"not null"`
	UpdatedAt time.Time `json:"-" gorm:"not null"`
}
