package validators

import (
	"time"

	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/sanctumlabs/curtz/app/pkg"
	"github.com/sanctumlabs/curtz/app/pkg/errdefs"

	validation "github.com/go-ozzo/ozzo-validation"
)

// IsValidUrl validates a url
func IsValidUrl(url string) error {
	err := validation.Validate(url, validation.Required, is.URL)
	if err != nil {
		return errdefs.ErrInvalidURL
	}
	return nil
}

// IsValidUserId checks if a given userId is valid
func IsValidUserId(userId string) error {
	if userId == "" {
		return errdefs.ErrInvalidUserId
	}
	return nil
}

// IsValidUrlId checks if a given urlId is valid
func IsValidUrlId(urlId string) error {
	if urlId == "" {
		return errdefs.ErrInvalidUrlId
	}
	return nil
}

// IsValidCustomAlias checks if a given custom alias is valid
func IsValidCustomAlias(customAlias string) error {
	if customAlias != "" && len(customAlias) > pkg.ShortCodeLength {
		return errdefs.ErrInvalidCustomAlias
	}
	return nil
}

// IsValidExpirationTime checks if the given timestamp is a valid expiration timestamp
func IsValidExpirationTime(expirationTime time.Time) error {
	if expirationTime.In(time.UTC).Before(time.Now().In(time.UTC)) {
		return errdefs.ErrPastExpiration
	}
	return nil
}
