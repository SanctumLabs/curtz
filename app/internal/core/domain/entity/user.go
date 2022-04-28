package entity

import (
	"github.com/sanctumlabs/curtz/app/internal/core/entity"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

const (
	// USER_STATUS_ACTIVE UserStatusActive is the status for an active user
	USER_STATUS_ACTIVE      = "ACTIVE"
	USER_STATUS_DEACTIVATED = "DEACTIVATED"
)

type User struct {
	identifier.ID
	entity.BaseEntity
	Email
	Password
	Token
	Verified bool
	Status   string
}

func NewUser(email, password string) (User, error) {
	userPassword := NewPassword(password)
	userEmail := NewEmail(email)
	id := identifier.New[User]()

	if !userEmail.isValid() {
		panic("Invalid email")
	}

	if err := userPassword.HashPassword(); err != nil {
		return User{}, err
	}

	userToken := NewToken()
	baseModel := entity.NewBaseEntity()

	return User{
		ID:         id,
		Email:      userEmail,
		Password:   userPassword,
		Token:      userToken,
		BaseEntity: baseModel,
	}, nil
}

func (user User) Prefix() string {
	return "user"
}
