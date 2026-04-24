package errdefs

import "errors"

var (
	ErrInvalidUserId                = errors.New("user id is invalid")
	ErrInvalidUrlId                 = errors.New("url id is invalid")
	ErrInvalidCustomAlias           = errors.New("custom alias is invalid")
	ErrServerError                  = errors.New("unexpected error encountered in server side")
	ErrInvalidURL                   = errors.New("url is invalid")
	ErrFilteredURL                  = errors.New("url matches filter pattern")
	ErrURLAlreadyExists             = errors.New("url matches filter pattern")
	ErrURLNotFound                  = errors.New("url not found")
	ErrURLExpired                   = errors.New("url expired")
	ErrEmailInvalid                 = errors.New("email is invalid")
	ErrInvalidPasswordLen           = errors.New("password length is invalid")
	ErrPasswordMissmatch            = errors.New("passwords don't match")
	ErrInvalidDate                  = errors.New("expires_on should be in 'yyyy-mm-dd hh:mm:ss' format")
	ErrURLAlreadyShort              = errors.New("url is already shortened")
	ErrNoMatchingData               = errors.New("no data matching given criteria")
	ErrShortCodeEmpty               = errors.New("short_code must not be empty")
	ErrNoShortCode                  = errors.New("short_code is not found")
	ErrTokenRequired                = errors.New("auth token is required")
	ErrTokenInvalid                 = errors.New("auth token is invalid")
	ErrUserExists                   = errors.New("User already exists")
	ErrUserDoestNotExist            = errors.New("User does not exist")
	ErrCustomAliasInvalidLength     = errors.New("custom alias must be between 3 and 100 characters")
	ErrCustomAliasInvalidCharacters = errors.New("custom alias can only contain alphanumeric characters and dashes")
	ErrInvalidURLLen                = errors.New("original URL must be between 5 and 2048 characters")
	ErrPastExpiration               = errors.New("expiration date cannot be in the past")
	ErrKeywordsCount                = errors.New("cannot have more than 10 keywords")
	ErrKeywordLength                = errors.New("keyword must be between 2 and 25 characters")
	ErrInvalidKeyword               = errors.New("keyword can only contain alphanumeric characters, dashes, and underscores")
	ErrShortCodeInvalidLength       = errors.New("short code must be between 6 and 10 characters")
	ErrShortCodeInvalidCharacters   = errors.New("short code can only contain alphanumeric characters")
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
