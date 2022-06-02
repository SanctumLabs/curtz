package validators

import (
	"github.com/go-ozzo/ozzo-validation/is"
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
