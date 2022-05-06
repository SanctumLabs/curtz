package entity

import (
	"regexp"
	"strings"
)

type Email struct {
	Value string
}

func NewEmail(value string) Email {
	return Email{
		Value: value,
	}
}

func (e Email) isValid() bool {
	pattern := `^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`
	re := regexp.MustCompile(pattern)
	return re.Match([]byte(e.Value))
}

func (e Email) format() string {
	return strings.ToLower(e.Value)
}
