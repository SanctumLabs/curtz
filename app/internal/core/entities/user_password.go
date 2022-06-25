package entities

import "github.com/sanctumlabs/curtz/app/pkg/utils"

// Password contains user password information
type Password struct {
	Value string
}

// NewPassword creates a new password
func NewPassword(value string) (Password, error) {
	hash, err := utils.HashPassword(value)

	if err != nil {
		return Password{}, err
	}

	return Password{
		Value: hash,
	}, nil
}

//Compare compares the password hash against the passed in password string
func (p Password) Compare(hash, password string) (bool, error) {
	ok, err := utils.CompareHashAndPassword(hash, password)
	if err != nil {
		return false, err
	}
	return ok, nil
}
