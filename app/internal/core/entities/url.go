package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/sanctumlabs/curtz/app/pkg/encoding"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
	"github.com/sanctumlabs/curtz/app/pkg/validators"
)

// URL represents an entity for a url
type URL struct {
	// ID is the unique identifier for the url
	identifier.ID

	// UserID is the user id of the url owner
	UserId identifier.ID

	// shortCode is the short code for the url
	shortCode string

	// CustomAlias is the custom alias for the url
	customAlias string

	// OriginalURL is the original url
	originalUrl string

	// Hits is the number of hits for the url
	hits uint

	// ExpiresOn is the expiration date for the url
	expiresOn time.Time

	// Keywords is a list of keywords for the url
	keywords []Keyword

	// BaseEntity is the base entity for the url
	BaseEntity
}

// NewUrl creates a new URL entity
func NewUrl(userId identifier.ID, originalUrl string, customAlias string, expiresOn time.Time, keywords []string) (*URL, error) {
	err := validators.IsValidUrl(originalUrl)
	if err != nil {
		return nil, err
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
		originalUrl: originalUrl,
		shortCode:   shortCode,
		customAlias: customAlias,
		keywords:    kws,
		expiresOn:   expiresOn,
	}, nil
}

// GetShortCode gets the short code of a given URL
func (url *URL) GetShortCode() string {
	return url.shortCode
}

// SetShortCode sets the short code on a URL
func (url *URL) SetShortCode(shortCode string) error {
	err := validators.IsValidShortCode(shortCode)
	if err != nil {
		return err
	}
	url.shortCode = shortCode
	return nil
}

// GetCustomAlias returns the custom alias for a URL
func (url *URL) GetCustomAlias() string {
	return url.customAlias
}

// SetCustomAlias sets the custom alias for a URL
func (url *URL) SetCustomAlias(customAlias string) error {
	err := validators.IsValidCustomAlias(customAlias)
	if err != nil {
		return err
	}
	url.customAlias = customAlias
	return nil
}

// GetOriginalURL returns the original url of URL
func (url *URL) GetOriginalURL() string {
	return url.originalUrl
}

// SetOriginalURL sets/update an original URL
func (url *URL) SetOriginalURL(u string) error {
	err := validators.IsValidUrl(u)
	if err != nil {
		return err
	}
	return nil
}

// GetHits returns the number of times the URL has been visited
func (url *URL) GetHits() uint {
	return url.hits
}

// SetHits sets/updates the number of times a URL has been visited
func (url *URL) SetHits(hits int) error {
	if hits < 0 {
		return errors.New("URL Visit count/hits can not be set to less than 0")
	}

	url.hits = uint(hits)
	return nil
}

// GetExpiresOn returns the expiry time of the URL
func (url *URL) GetExpiresOn() time.Time {
	return url.expiresOn
}

// SetExpiresOn updates/sets the expiry time of the URL
func (url *URL) SetExpiresOn(expiresOn time.Time) {
	url.expiresOn = expiresOn
}

// IsActive checks if the url is active or not expired.
func (url *URL) IsActive() bool {
	return url.expiresOn.In(time.UTC).After(time.Now().In(time.UTC))
}

// GetExpiryDuration returns as a time.Duration how long before the url expires
// This returns an absolute value after subtracting time.Now()
func (url *URL) GetExpiryDuration() time.Duration {
	duration := time.Until(url.expiresOn)
	if duration >= 0 {
		return duration
	}
	return -duration
}

// GetKeywords returns the keywords on the url
func (url *URL) GetKeywords() []Keyword {
	return url.keywords
}

// SetKeywords adds new keywords to the url
func (url *URL) SetKeywords(keywords []string) error {
	kws, err := createKeywords(keywords)
	if err != nil {
		return err
	}

	url.keywords = append(url.keywords, kws...)
	return nil
}

// Prefix returns the url prefix for logging
func (url *URL) Prefix() string {
	return fmt.Sprintf("url-%s-%s", url.ID, url.shortCode)
}
