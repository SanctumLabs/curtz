package contracts

import (
	"github.com/google/uuid"
	"github.com/sanctumlabs/curtz/app/internal/core/domain/entity"
)

type UrlRepository interface {
	Save(owner uuid.UUID, originalUrl, shortenedUrl string) (entity.URL, error)
	GetByShortUrl(shortenedUrl string) (entity.URL, error)
	GetByOwner(owner uuid.UUID) ([]entity.URL, error)
	GetByKeyword(keyword string) ([]entity.URL, error)
	GetByKeywords(keywords []string) ([]entity.URL, error)
	GetByOriginalUrl(originalUrl string) ([]entity.URL, error)
	GetById(id uuid.UUID) (entity.URL, error)
	Delete(id uuid.UUID) error
}

type UserRepository interface {
	CreateUser(email, password string) (uuid.UUID, error)
	GetByEmail(email string) (uuid.UUID, error)
	GetById(id uuid.UUID) (uuid.UUID, error)
}
