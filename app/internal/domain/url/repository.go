package url

import (
	"time"
)

type UrlRepository interface {
	UrlReadRepository
	UrlWriteRepository
}

type UrlWriteRepository interface {
	Save(URL) (URL, error)
	Update(urlID, customAlias string, keywords []Keyword, expiresOn *time.Time) (URL, error)
	Delete(id string) error
	IncrementHits(shortCode string) error
}

type UrlReadRepository interface {
	GetByShortCode(shortCode string) (URL, error)
	GetByOwner(owner string) ([]URL, error)
	GetByKeyword(keyword string) ([]URL, error)
	GetByKeywords(keywords []string) ([]URL, error)
	GetByOriginalUrl(originalUrl string) (URL, error)
	GetById(id string) (URL, error)
}
