package url

import (
	"time"

	"github.com/google/uuid"
	"github.com/sanctumlabs/curtz/internal/core/domain/entities"
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

// IsActive checks if the url model is active
// It returns true if url is not marked deleted or expired, false otherwise.
func (url URL) IsActive() bool {
	if url.Deleted {
		return false
	}

	return url.ExpiresOn.In(time.UTC).After(time.Now().In(time.UTC))
}
