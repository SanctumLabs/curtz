package errdefs

import "errors"

var (
	ErrInvalidUserId      = errors.New("user id is invalid")
	ErrURLIdInvalid       = errors.New("url id is invalid")
	ErrInvalidCustomAlias = errors.New("custom alias is invalid")
	ErrServerError        = errors.New("unexpected error encountered in server side")
	ErrURLInvalid         = errors.New("url is invalid")
	ErrURLLength          = errors.New("url is too short or too long, should be 15-2048 chars")
	ErrURLFiltered        = errors.New("url matches filter pattern")
	ErrURLAlreadyExists   = errors.New("url already exists")
	ErrURLNotFound        = errors.New("url not found")
	ErrURLExpired         = errors.New("url expired")
	ErrURLAlreadyShort    = errors.New("url is already shortened")
	ErrKeywordsCount      = errors.New("keywords must not be more than 10")
	ErrKeywordLength      = errors.New("keyword must contain 2-25 characters")
	ErrInvalidKeyword     = errors.New("keyword must be alphanumeric (dash/underscore allowed)")
	ErrEmailInvalid       = errors.New("email is invalid")
	ErrInvalidPasswordLen = errors.New("password length is invalid")
	ErrPasswordMissmatch  = errors.New("passwords don't match")
	ErrInvalidDate        = errors.New("expires_on should be in 'yyyy-mm-dd hh:mm:ss' format")
	ErrPastExpiration     = errors.New("expires_on can not be a date in the past")
	ErrNoMatchingData     = errors.New("no data matching given criteria")
	ErrShortCodeEmpty     = errors.New("short code must not be empty")
	ErrNoShortCode        = errors.New("short code is not found")
	ErrShortCodeInvalid   = errors.New("short code is invalid")
	ErrTokenRequired      = errors.New("auth token is required")
	ErrTokenInvalid       = errors.New("auth token is invalid")
	ErrUserExists         = errors.New("User already exists")
	ErrUserDoestNotExist  = errors.New("User does not exist")
)

type Error struct {
	msg string
}

func (e *Error) Error() string {
	return e.msg
}

func NewError(msg string) Error {
	return Error{msg}
}
