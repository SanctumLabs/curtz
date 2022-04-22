package user

import "errors"

var (
	ErrTokenRequired = errors.New("auth token is required")
	ErrTokenInvalid  = errors.New("auth token is invalid")
)
