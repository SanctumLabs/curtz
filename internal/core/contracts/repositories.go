package contracts

import (
	"github.com/google/uuid"
	"github.com/sanctumlabs/curtz/internal/core/domain/url"
)

type UrlRepository interface {
	CreateUrl(owner uuid.UUID, originalUrl, shortenedUrl string) (url.URL, error)
	GetByShortUrl(shortenedUrl string) (url.URL, error)
	GetByOwner(owner uuid.UUID) ([]url.URL, error)
	GetByKeyword(keyword string) ([]url.URL, error)
	GetByKeywords(keywords []string) ([]url.URL, error)
	GetByOriginalUrl(originalUrl string) ([]url.URL, error)
	GetById(id uuid.UUID) (url.URL, error)
	Delete(id uuid.UUID) error
}

type UserRepository interface {
	CreateUser(email, password string) (uuid.UUID, error)
	GetByEmail(email string) (uuid.UUID, error)
	GetById(id uuid.UUID) (uuid.UUID, error)
}
