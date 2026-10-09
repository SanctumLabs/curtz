//go:build legacy

package entities

import (
	"github.com/sanctumlabs/fupi/app/pkg/errdefs"
	"github.com/sanctumlabs/fupi/app/pkg/utils"
)

// Email is a domain entity for email
type Email struct {
	Value string
}

// NewEmail creates a new email
func NewEmail(value string) (Email, error) {
	if utils.IsEmailValid(value) {
		return Email{
			Value: value,
		}, nil
	}
	return Email{}, errdefs.ErrEmailInvalid
}
