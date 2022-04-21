package user

import (
	"github.com/sanctumlabs/curtz/internal/core/domain/entities"
)

const (
	// USER_STATUS_ACTIVE UserStatusActive is the status for an active user
	USER_STATUS_ACTIVE      = "ACTIVE"
	USER_STATUS_DEACTIVATED = "DEACTIVATED"
)

type User struct {
	entities.Identifier
	entities.BaseEntity
	Email
	Password
	Token
	Verified bool   `db:"verified,omitempty"`
	Status   string `db:"status,omitempty"`
}

func NewUser(email, password string) (User, error) {
	userPassword := NewPassword(password)
	userEmail := NewEmail(email)
	identifier := entities.NewIdentifier()

	if !userEmail.isValid() {
		panic("Invalid email")
	}

	if err := userPassword.HashPassword(); err != nil {
		return User{}, err
	}

	userToken := NewToken()
	baseModel := entities.NewBaseEntity()

	return User{
		Identifier: identifier,
		Email:      userEmail,
		Password:   userPassword,
		Token:      userToken,
		BaseEntity: baseModel,
	}, nil
}
