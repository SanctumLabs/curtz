package entities

import "time"

// BaseEntity is a base model for all models
type BaseEntity struct {
	// DeletedAt is the date when the entity was deleted
	DeletedAt *time.Time

	// CreatedAt is the date when the entity was created
	CreatedAt time.Time

	// UpdatedAt is the date when the entity was updated
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
