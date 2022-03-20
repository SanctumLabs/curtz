package models

import (
	"github.com/google/uuid"
	"github.com/sanctumlabs/curtz/pkg/uidgen"
)

// Identifier is a model for an Identifier
type Identifier struct {
	ID   uint      `json:"-" gorm:"primaryKey"`
	UUID uuid.UUID `json:"id" gorm:"type:uuid;uniqueIndex;not null"`
}

func NewIdentifier() Identifier {
	uuid := uidgen.New().Create()
	return Identifier{
		UUID: uuid,
	}
}
