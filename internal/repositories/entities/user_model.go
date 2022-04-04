package entities

import (
	"github.com/google/uuid"
	"time"
)

type UserModel struct {
	BaseModel
	Email                string     `gorm:"column:email,uniqueIndex"`
	Password             string     `gorm:"column:password"`
	APIKey               *uuid.UUID `gorm:"column:api_key"`
	ResetPasswordExpires *time.Time `gorm:"column:reset_password_expires"`
	ResetPasswordToken   *uuid.UUID `gorm:"column:reset_password_token"`
	VerificationExpires  time.Time  `gorm:"column:verification_expires"`
	VerificationToken    uuid.UUID  `gorm:"column:verification_token"`
	Verified             bool       `gorm:"column:verified,omitempty"`
	Urls                 []UrlModel `gorm:"foreignkey:UserID"`
}
