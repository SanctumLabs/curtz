package models

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel definition same as gorm.Model, but including other common columns
type BaseModel struct {
	gorm.Model
	ID         uint       `cql:"id" gorm:"primary_key;column:id"`
	Identifier uint       `cql:"identifier" gorm:"column;column:identifier"`
	CreatedAt  time.Time  `cql:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time  `cql:"updated_at" gorm:"column:updated_at"`
	DeletedAt  *time.Time `cql:"deleted_at" gorm:"column:deleted_at"`
}
