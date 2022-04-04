package contracts

import (
	"github.com/google/uuid"
	"github.com/sanctumlabs/curtz/internal/core/domain/url"
)

type UrlService interface {
	CreateUrl(owner uuid.UUID, originalUrl, shortenedUrl string) (url.URL, error)
	GetByShortUrl(shortenedUrl string) (url.URL, error)
	GetByOwner(owner uuid.UUID) ([]url.URL, error)
	GetByKeyword(keyword string) ([]url.URL, error)
	GetByKeywords(keywords []string) ([]url.URL, error)
	GetByOriginalUrl(originalUrl string) ([]url.URL, error)
	GetById(id uuid.UUID) (url.URL, error)
	Remove(id uuid.UUID) error
}
