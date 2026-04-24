package url

import (
	"fmt"
	netUrl "net/url"
	"time"

	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/pkg/encoding"
	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

// URL represents an entity for a url
type (
	URL struct {
		// ID is the unique identifier for the url
		identifier.ID

		// UserID is the user id of the url owner
		userId identifier.ID

		// ShortCode is the short code for the url
		shortCode ShortCode

		// CustomAlias is the custom alias for the url
		customAlias CustomAlias

		// OriginalURL is the original url
		originalUrl OriginalURL

		// Keywords is a list of keywords for the url
		Keywords []Keyword

		status URLStatus

		// ExpiresOn is the expiration date for the url
		ExpiresOn time.Time

		// BaseEntity is the base entity for the url
		entities.BaseEntity
	}

	// URLParams represents the parameters for creating or updating a url
	URLParams struct {
		// ID is the unique identifier for the url
		ID string

		// UserID is the user id of the url owner
		UserId string

		// ShortCode is the short code for the url
		ShortCode string

		// CustomAlias is the custom alias for the url
		CustomAlias string

		// OriginalURL is the original url
		OriginalUrl string

		// Hits is the number of hits for the url
		Hits uint

		// ExpiresOn is the expiration date for the url
		ExpiresOn time.Time

		// Keywords is a list of keywords for the url
		Keywords []string

		Status URLStatus
	}
)

// NewUrl creates a new URL entity
func NewUrl(params URLParams) (*URL, error) {
	if l := len(originalUrl); l < MinLength || l > MaxLength {
		return nil, errdefs.ErrInvalidURLLen
	}

	if filterRe.MatchString(originalUrl) {
		return nil, errdefs.ErrFilteredURL
	}

	_, err := netUrl.ParseRequestURI(originalUrl)
	if err != nil {
		return nil, errdefs.ErrInvalidURL
	}

	if !urlRe.MatchString(originalUrl) {
		return nil, errdefs.ErrInvalidURL
	}

	if expiresOn.In(time.UTC).Before(time.Now().In(time.UTC)) {
		return nil, errdefs.ErrPastExpiration
	}

	id := identifier.New()

	shortCode, err := encoding.GetUniqueShortCode()
	if err != nil {
		return nil, err
	}

	kws, err := createKeywords(keywords)
	if err != nil {
		return nil, err
	}

	return &URL{
		ID:          id,
		UserId:      userId,
		BaseEntity:  NewBaseEntity(),
		OriginalUrl: originalUrl,
		ShortCode:   shortCode,
		CustomAlias: customAlias,
		Keywords:    kws,
		ExpiresOn:   expiresOn,
	}, nil
}

// IsActive checks if the url is active or not expired.
func (url URL) IsActive() bool {
	return url.ExpiresOn.In(time.UTC).After(time.Now().In(time.UTC))
}

func (url URL) GetKeywords() []Keyword {
	return url.Keywords
}

func (url URL) SetKeywords(keywords []string) error {
	kws, err := createKeywords(keywords)
	if err != nil {
		return err
	}

	url.Keywords = append(url.Keywords, kws...)
	return nil
}

func (url URL) GetExpiresOn() time.Time {
	return url.ExpiresOn
}

func (url URL) SetExpiresOn(expiresOn time.Time) {
	url.ExpiresOn = expiresOn
}

// UpdateExpiresOn updates the expiresOn field on a URL
func (url URL) UpdateExpiresOn(expiresOn time.Time) error {
	url.ExpiresOn = expiresOn
	return nil
}

func (url URL) GetCustomAlias() string {
	return url.CustomAlias
}

func (url URL) SetCustomAlias(customAlias string) error {
	url.CustomAlias = customAlias
	return nil
}

// GetExpiryDuration returns as a time.Duration how long before the url expires
// This returns an absolute value after subtracting time.Now()
func (url URL) GetExpiryDuration() time.Duration {
	duration := time.Until(url.ExpiresOn)
	if duration >= 0 {
		return duration
	}
	return -duration
}

// Prefix returns the url prefix for logging
func (url URL) Prefix() string {
	return fmt.Sprintf("url-%s-%s", url.ID, url.ShortCode)
}
