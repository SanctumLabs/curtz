package user

import (
	"time"

	"github.com/google/uuid"
)

// UserToken contains token information for a user
type UserToken struct {
	APIKey               *uuid.UUID `db:"api_key,omitempty"`
	ResetPasswordExpires *time.Time `db:"reset_password_expires,omitempty"`
	ResetPasswordToken   *uuid.UUID `db:"reset_password_token,omitempty"`
	VerificationExpires  time.Time  `db:"verification_expires"`
	VerificationToken    uuid.UUID  `db:"verification_token,omitempty"`
}
