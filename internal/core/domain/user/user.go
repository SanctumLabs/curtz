package user

import (
	"time"

	"github.com/sanctumlabs/curtz/internal/core/domain/models"
	"github.com/sanctumlabs/curtz/pkg/encoding"
)

const (
	// UserStatusActive is the status for an active user
	USER_STATUS_ACTIVE      = "ACTIVE"
	USER_STATUS_DEACTIVATED = "DEACTIVATED"
)

type User struct {
	models.Identifier
	models.BaseModel
	UserEmail
	UserPassword
	UserToken
	Verified bool   `db:"verified,omitempty"`
	Status   string `db:"status,omitempty"`
}

func NewUser(email, password string) User {
	var now time.Time = time.Now()
	userPassword := NewUserPassword(password)
	userEmail := NewUserEmail(email)
	identifier := models.NewIdentifier()

	if !userEmail.isValid() {
		panic("Invalid email")
	}

	userPassword.HashPassword()

	verificationToken := encoding.GenUniqueID()
	verificationToken.String()

	return User{
		Identifier:   identifier,
		UserEmail:    userEmail,
		UserPassword: userPassword,
		UserToken: {
			VerificationToken:   verificationToken.String(),
			VerificationExpires: now.Add(time.Minute * 15).UTC().Round(time.Microsecond), //expires 15 mins later
		},
		BaseModel: {
			CreatedAt: now.UTC().Round(time.Microsecond),
			UpdatedAt: now.UTC().Round(time.Microsecond),
		},
	}
}
