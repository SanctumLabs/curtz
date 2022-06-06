package entity

import (
	"fmt"
	netUrl "net/url"
	"regexp"
	"time"

	"github.com/sanctumlabs/curtz/app/internal/core/entity"
	"github.com/sanctumlabs/curtz/app/pkg"
	"github.com/sanctumlabs/curtz/app/pkg/encoding"
	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

// @see https://github.com/asaskevich/govalidator/blob/master/patterns.go
var (
	IP        = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	Schema    = `((ftp|https?):\/\/)`
	Username  = `(\S+(:\S*)?@)`
	Path      = `((\/|\?|#)[^\s]*)`
	Port      = `(:(\d{1,5}))`
	UrlIP     = `([1-9]\d?|1\d\d|2[01]\d|22[0-3]|24\d|25[0-5])(\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-5]))`
	Subdomain = `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`

	MinLength   = 15
	MaxLength   = 2048
	Regex       = `^` + Schema + `?` + Username + `?` + `((` + IP + `|(\[` + IP + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + Subdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + Port + `?` + Path + `?$`
	FilterRegex = `(xxx|localhost|127\.0\.0\.1|\.lvh\.me|\.local|urlss?h\.)`
)

var (
	urlRe    = regexp.MustCompile(Regex)
	filterRe = regexp.MustCompile(FilterRegex)
)

// URL represents an entity for a url
type URL struct {
	// ID is the unique identifier for the url
	identifier.ID

	// UserID is the user id of the url owner
	UserId identifier.ID

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
	Keywords []Keyword

	// BaseEntity is the base entity for the url
	entity.BaseEntity
}

// NewUrl creates a new URL entity
func NewUrl(userId identifier.ID, originalUrl, customAlias, expiresOn string, keywords []string) (*URL, error) {
	kws := make([]Keyword, len(keywords))

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

	if len(expiresOn) != len(pkg.DateLayout) {
		return nil, errdefs.ErrInvalidDate
	}

	expiry, err := parseExpiresOn(expiresOn)

	if err != nil {
		return nil, errdefs.ErrInvalidDate
	}

	if expiry.In(time.UTC).Before(time.Now().In(time.UTC)) {
		return nil, errdefs.ErrPastExpiration
	}

	id := identifier.New()

	shortCode, err := encoding.GetUniqueShortCode()
	if err != nil {
		return nil, err
	}

	kws, err = createKeywords(keywords)
	if err != nil {
		return nil, err
	}

	return &URL{
		ID:          id,
		UserId:      userId,
		BaseEntity:  entity.NewBaseEntity(),
		OriginalUrl: originalUrl,
		ShortCode:   shortCode,
		CustomAlias: customAlias,
		Keywords:    kws,
		ExpiresOn:   expiry,
	}, nil
}

// IsActive checks if the url model is active
// It returns true if url is not marked deleted or expired, false otherwise.
func (url URL) IsActive() bool {
	return url.ExpiresOn.In(time.UTC).After(time.Now().In(time.UTC))
}

func parseExpiresOn(expiresOn string) (time.Time, error) {
	if expiresOn == "" {
		return time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC), nil
	}

	return time.ParseInLocation(pkg.DateLayout, expiresOn, time.UTC)
}

// Prefix returns the url prefix for logging
func (url URL) Prefix() string {
	return fmt.Sprintf("url-%s-%s", url.ID, url.ShortCode)
}
