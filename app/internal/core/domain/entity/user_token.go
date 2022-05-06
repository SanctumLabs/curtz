package entity

import (
	"time"

	"github.com/sanctumlabs/curtz/app/pkg/encoding"

	"github.com/google/uuid"
)

// Token contains token information for a user
type Token struct {
	APIKey               *uuid.UUID
	ResetPasswordExpires *time.Time
	ResetPasswordToken   *uuid.UUID
	VerificationExpires  time.Time
	VerificationToken    uuid.UUID
}

func NewToken() Token {
	var now time.Time = time.Now()
	verificationToken := encoding.GenUniqueID()
	return Token{
		VerificationToken:   verificationToken,
		VerificationExpires: now.Add(time.Minute * 15).UTC().Round(time.Microsecond), //expires 15 mins later
	}
}
