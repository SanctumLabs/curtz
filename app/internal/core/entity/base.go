package entity

import "time"

// BaseEntity is a base model for all models
type BaseEntity struct {
	DeletedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewBaseEntity creates a new base entity
func NewBaseEntity() BaseEntity {
	now := time.Now()
	return BaseEntity{
		CreatedAt: now.UTC().Round(time.Microsecond),
		UpdatedAt: now.UTC().Round(time.Microsecond),
	}
}
