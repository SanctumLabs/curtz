package entities

import (
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

const (
	// USER_STATUS_ACTIVE UserStatusActive is the status for an active user
	USER_STATUS_ACTIVE      = "ACTIVE"
	USER_STATUS_DEACTIVATED = "DEACTIVATED"
)

// User is a domain entity for user
type User struct {
	// ID is the unique identifier for a user
	id identifier.ID

	// BaseEntity is the base entity for a user
	BaseEntity

	// Email is the email address for a user
	email *Email

	// Password is the password for a user
	password Password

	// Token is the token for a user
	Token

	// Verified is the verification status for a user
	Verified bool

	// Status is the status for a user
	Status string
}

func NewUser(email, password string) (*User, error) {
	userPassword, err := NewPassword(password)
	if err != nil {
		return nil, err
	}

	userEmail, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	id := identifier.New()
	userToken := NewToken()
	baseModel := NewBaseEntity()

	return &User{
		id:         id,
		email:      userEmail,
		password:   userPassword,
		Token:      userToken,
		BaseEntity: baseModel,
	}, nil
}

// GetId retrieves the user id
func (u *User) GetId() string {
	return u.id.String()
}

func (u *User) SetId(id string) error {
	uid, err := identifier.New().FromString(id)
	if err != nil {
		return err
	}
	u.id = uid
	return nil
}

// GetEmail retrieves the user email
func (u *User) GetEmail() string {
	return u.email.GetValue()
}

// SetEmail sets the email of the user
func (u *User) SetEmail(email string) error {
	newEmail, err := NewEmail(email)
	if err != nil {
		return err
	}
	u.email = newEmail
	return nil
}

// GetPassword returns the user password
func (u *User) GetPassword() string {
	return u.password.GetValue()
}

// CheckPassword compares the hash of the password and the password
func (u *User) CheckPassword(password string) (bool, error) {
	hash := u.password.value
	return u.password.Compare(hash, password)
}

// SetPassword sets the user password
func (u *User) SetPassword(password string) error {
	newPassword, err := NewPassword(password)
	if err != nil {
		return err
	}
	u.password = newPassword
	return nil
}
