package models

import (
	"time"

	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

type User struct {
	BaseModel
	Email                string         `gorm:"column:email,uniqueIndex"`
	Password             string         `gorm:"column:password"`
	ResetPasswordExpires *time.Time     `gorm:"column:reset_password_expires"`
	ResetPasswordToken   *identifier.ID `gorm:"column:reset_password_token"`
	VerificationExpires  time.Time      `gorm:"column:verification_expires"`
	VerificationToken    identifier.ID  `gorm:"column:verification_token"`
	Verified             bool           `gorm:"column:verified,omitempty"`
	Urls                 []Url          `gorm:"foreignKey:Owner"`
}
