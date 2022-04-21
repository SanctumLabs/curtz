package entities

import "time"

// BaseEntity is a base model for all models
type BaseEntity struct {
	Deleted   bool      `json:"is_deleted" gorm:"default:false;not null"`
	CreatedAt time.Time `json:"-" gorm:"not null"`
	UpdatedAt time.Time `json:"-" gorm:"not null"`
}

func NewBaseEntity() BaseEntity {
	now := time.Now()
	return BaseEntity{
		CreatedAt: now.UTC().Round(time.Microsecond),
		UpdatedAt: now.UTC().Round(time.Microsecond),
	}
}
