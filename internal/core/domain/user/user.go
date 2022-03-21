package user

import (
	"github.com/sanctumlabs/curtz/internal/core/domain/models"
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

func NewUser(email, password string) (User, error) {
	userPassword := NewUserPassword(password)
	userEmail := NewUserEmail(email)
	identifier := models.NewIdentifier()

	if !userEmail.isValid() {
		panic("Invalid email")
	}

	if err := userPassword.HashPassword(); err != nil {
		return User{}, err
	}

	userToken := NewUserToken()
	baseModel := models.NewBaseModel()

	return User{
		Identifier:   identifier,
		UserEmail:    userEmail,
		UserPassword: userPassword,
		UserToken:    userToken,
		BaseModel:    baseModel,
	}, nil
}
