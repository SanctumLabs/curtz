package url

import (
	"fmt"
	"time"

	"github.com/sanctumlabs/curtz/app/internal/core/entity"
	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
)

type (
	// URL is the aggregate root for the URL bounded context.
	// It owns all state transitions for a shortened link.
	URL struct {
		entity.AggregateRoot

		// UserID is the user id of the url owner
		userId entity.ID

		// ShortCode is the short code for the url
		shortCode ShortCode

		// CustomAlias is the custom alias for the url
		customAlias CustomAlias

		// OriginalURL is the original url
		originalUrl OriginalURL

		// Keywords is a list of keywords for the url
		keywords []Keyword

		status URLStatus

		// ExpiresOn is the expiration date for the url
		expiresOn time.Time
	}

	// URLParams represents the parameters for creating or updating a url
	URLParams struct {
		entity.AggregateRootParams

		// UserID is the user id of the url owner
		UserId string

		// ShortCode is the short code for the url
		ShortCode string

		// CustomAlias is the custom alias for the url
		CustomAlias string

		// OriginalURL is the original url
		OriginalUrl string

		// ExpiresOn is the expiration date for the url
		ExpiresOn time.Time

		// Keywords is a list of keywords for the url
		Keywords []string

		Status URLStatus
	}
)

// NewUrl creates a new URL entity
func NewUrl(params URLParams) (*URL, error) {
	originalUrl, err := NewOriginalURL(params.OriginalUrl)
	if err != nil {
		return nil, err
	}

	if params.ExpiresOn.In(time.UTC).Before(time.Now().In(time.UTC)) {
		return nil, errdefs.ErrPastExpiration
	}

	kws, err := createKeywords(params.Keywords)
	if err != nil {
		return nil, err
	}

	customAlias, customAliasErr := NewCustomAlias(params.CustomAlias)
	if customAliasErr != nil {
		return nil, customAliasErr
	}

	shortCode, shortCodeErr := NewShortCode(params.ShortCode)
	if shortCodeErr != nil {
		return nil, shortCodeErr
	}

	aggregateRoot, err := entity.NewAggregateRoot(params.AggregateRootParams)
	if err != nil {
		return nil, err
	}

	userId, userIdErr := entity.StringToID(params.UserId)
	if userIdErr != nil {
		return nil, userIdErr
	}

	return &URL{
		AggregateRoot: aggregateRoot,
		shortCode:     shortCode,
		customAlias:   customAlias,
		userId:        userId,
		originalUrl:   originalUrl,
		keywords:      kws,
		expiresOn:     params.ExpiresOn,
		status:        params.Status,
	}, nil
}

// IsActive checks if the url is active or not expired.
func (url *URL) IsActive() bool {
	return url.expiresOn.In(time.UTC).After(time.Now().In(time.UTC))
}

func (url *URL) OriginalURL() OriginalURL {
	return url.originalUrl
}

func (url *URL) ShortCode() ShortCode {
	return url.shortCode
}

func (url *URL) Keywords() []Keyword {
	return url.keywords
}

func (url *URL) ExpiresOn() time.Time {
	return url.expiresOn
}

func (url *URL) CustomAlias() CustomAlias {
	return url.customAlias
}

func (url *URL) Status() URLStatus {
	return url.status
}

// ExpiryDuration returns as a time.Duration how long before the url expires
// This returns an absolute value after subtracting time.Now()
func (url *URL) ExpiryDuration() time.Duration {
	duration := time.Until(url.expiresOn)
	if duration >= 0 {
		return duration
	}
	return -duration
}

// Prefix returns the url prefix for logging
func (url *URL) Prefix() string {
	return fmt.Sprintf("url-%s-%s", url.ID(), url.shortCode)
}
