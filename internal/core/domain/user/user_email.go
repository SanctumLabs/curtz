package user

import (
	"regexp"
	"strings"
)

type Email struct {
	Email string `db:"email,omitempty"`
}

func NewEmail(value string) Email {
	return Email{
		Email: value,
	}
}

func (e Email) isValid() bool {
	pattern := `^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`
	re := regexp.MustCompile(pattern)
	return re.Match([]byte(e.Email))
}

func (e Email) format() string {
	return strings.ToLower(e.Email)
}
