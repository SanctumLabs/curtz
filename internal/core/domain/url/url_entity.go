package url

import (
	netUrl "net/url"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/sanctumlabs/curtz/internal/core/domain/entities"
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

// URL is model for short urls
type URL struct {
	entities.Identifier
	Owner        uuid.UUID `json:"owner" gorm:"owner_id"`
	ShortenedUrl string    `json:"short_code" gorm:"size:12;uniqueIndex;not null"`
	OriginalUrl  string    `json:"original_url" gorm:"size:2048;index;not null"`
	Hits         uint      `json:"hits" gorm:"default:0;not null"`
	entities.BaseEntity
	ExpiresOn time.Time `json:"expires_on"`
	Keywords  []Keyword `json:"-" gorm:"many2many:url_keywords"`
}

func NewUrl(owner uuid.UUID, originalUrl, shortenedUrl string) URL {
	identifier := entities.NewIdentifier()

	return URL{
		Identifier:   identifier,
		Owner:        owner,
		BaseEntity:   entities.NewBaseEntity(),
		OriginalUrl:  originalUrl,
		ShortenedUrl: shortenedUrl,
	}
}

// Validate validates the url
func (url *URL) Validate() error {
	if l := len(url.OriginalUrl); l < MinLength || l > MaxLength {
		return ErrInvalidURLLen
	}

	if filterRe.MatchString(url.OriginalUrl) {
		return ErrFilteredURL
	}

	uri, err := netUrl.ParseRequestURI(url.OriginalUrl)
	if err != nil {
		return ErrInvalidURL
	}

	if urlRe.MatchString(url.OriginalUrl) {
		return ErrInvalidURL
	}
}

// IsActive checks if the url model is active
// It returns true if url is not marked deleted or expired, false otherwise.
func (url URL) IsActive() bool {
	if url.Deleted {
		return false
	}

	return url.ExpiresOn.In(time.UTC).After(time.Now().In(time.UTC))
}
