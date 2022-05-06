package entity

import (
	netUrl "net/url"
	"regexp"
	"time"

	"github.com/sanctumlabs/curtz/app/internal/core/domain"
	"github.com/sanctumlabs/curtz/app/internal/core/entity"
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

	MinLength    = 15
	MaxLength    = 2048
	KeywordRegex = `^[a-zA-Z0-9-_]+$`
	Regex        = `^` + Schema + `?` + Username + `?` + `((` + IP + `|(\[` + IP + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + Subdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + Port + `?` + Path + `?$`
	FilterRegex  = `(xxx|localhost|127\.0\.0\.1|\.lvh\.me|\.local|urlss?h\.)`
)

var (
	kwRe     = regexp.MustCompile(KeywordRegex)
	urlRe    = regexp.MustCompile(Regex)
	filterRe = regexp.MustCompile(FilterRegex)
)

// URL is model for urls
type URL struct {
	identifier.ID
	UserId       identifier.ID
	ShortenedUrl string
	OriginalUrl  string
	Hits         uint
	entity.BaseEntity
	ExpiresOn time.Time
	Keywords  []Keyword
}

func New(userId identifier.ID, originalUrl string) (*URL, error) {
	id := identifier.New()

	if l := len(originalUrl); l < MinLength || l > MaxLength {
		return nil, domain.ErrInvalidURLLen
	}

	if filterRe.MatchString(originalUrl) {
		return nil, domain.ErrFilteredURL
	}

	_, err := netUrl.ParseRequestURI(originalUrl)
	if err != nil {
		return nil, domain.ErrInvalidURL
	}

	if urlRe.MatchString(originalUrl) {
		return nil, domain.ErrInvalidURL
	}

	// TODO: shorten url from long url
	shortenedUrl := ""

	return &URL{
		ID:           id,
		UserId:       userId,
		BaseEntity:   entity.NewBaseEntity(),
		OriginalUrl:  originalUrl,
		ShortenedUrl: shortenedUrl,
	}, nil
}

// IsActive checks if the url model is active
// It returns true if url is not marked deleted or expired, false otherwise.
func (url URL) IsActive() bool {
	if url.Deleted {
		return false
	}

	return url.ExpiresOn.In(time.UTC).After(time.Now().In(time.UTC))
}

func (url URL) Prefix() string {
	return "url"
}
