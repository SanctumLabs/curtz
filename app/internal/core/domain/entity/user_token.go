package entity

import (
	"github.com/sanctumlabs/curtz/app/pkg/encoding"
	"time"

	"github.com/google/uuid"
)

// Token contains token information for a user
type Token struct {
	APIKey               *uuid.UUID `db:"api_key,omitempty"`
	ResetPasswordExpires *time.Time `db:"reset_password_expires,omitempty"`
	ResetPasswordToken   *uuid.UUID `db:"reset_password_token,omitempty"`
	VerificationExpires  time.Time  `db:"verification_expires"`
	VerificationToken    uuid.UUID  `db:"verification_token,omitempty"`
}

func NewToken() Token {
	var now time.Time = time.Now()
	verificationToken := encoding.GenUniqueID()
	return Token{
		VerificationToken:   verificationToken,
		VerificationExpires: now.Add(time.Minute * 15).UTC().Round(time.Microsecond), //expires 15 mins later
	}
}
