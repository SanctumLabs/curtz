package user

import (
	"regexp"
	"strings"
)

type UserEmail struct {
	Email string `db:"email,omitempty"`
}

func NewUserEmail(value string) UserEmail {
	return UserEmail{
		Email: value,
	}
}

func (e UserEmail) isValid() bool {
	pattern := `^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`
	re := regexp.MustCompile(pattern)
	return re.Match([]byte(e.Email))
}

func (e UserEmail) format() string {
	return strings.ToLower(e.Email)
}
