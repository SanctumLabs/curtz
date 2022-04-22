package url

import "errors"

var (
	ErrInvalidURL      = errors.New("url is invalid")
	ErrInvalidURLLen   = errors.New("url is too short or too long, should be 15-2048 chars")
	ErrFilteredURL     = errors.New("url matches filter pattern")
	ErrKeywordsCount   = errors.New("keywords must not be more than 10")
	ErrKeywordLength   = errors.New("keyword must contain 2-25 characters")
	ErrInvalidKeyword  = errors.New("keyword must be alphanumeric (dash/underscore allowed)")
	ErrInvalidDate     = errors.New("expires_on should be in 'yyyy-mm-dd hh:mm:ss' format")
	ErrPastExpiration  = errors.New("expires_on can not be date in past")
	ErrURLAlreadyShort = errors.New("url is already shortened")
	ErrNoMatchingData  = errors.New("no data matching given criteria")
	ErrShortCodeEmpty  = errors.New("short_code must not be empty")
	ErrNoShortCode     = errors.New("short_code is not found")
)
