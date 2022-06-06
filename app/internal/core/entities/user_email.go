package entities

import (
	"regexp"
	"strings"

	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
)

// Email is a domain entity for email
type Email struct {
	Value string
}

// NewEmail creates a new email
func NewEmail(value string) (Email, error) {
	if isValid(value) {
		return Email{
			Value: value,
		}, nil
	}
	return Email{}, errdefs.ErrEmailInvalid
}

// isValid checks if an email address is valid
func isValid(value string) bool {
	pattern := `^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`
	re := regexp.MustCompile(pattern)
	return re.Match([]byte(value))
}

func (e Email) format() string {
	return strings.ToLower(e.Value)
}
