package entity

import "time"

// BaseEntity is a base model for all models
type BaseEntity struct {
	Deleted   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBaseEntity() BaseEntity {
	now := time.Now()
	return BaseEntity{
		CreatedAt: now.UTC().Round(time.Microsecond),
		UpdatedAt: now.UTC().Round(time.Microsecond),
	}
}
