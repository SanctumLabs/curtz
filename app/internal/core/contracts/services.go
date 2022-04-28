package contracts

import (
	"github.com/google/uuid"
	"github.com/sanctumlabs/curtz/app/internal/core/domain/entity"
)

type AuthService interface {
	// Authenticate a user given the token. Returns user id if authenticated, error otherwise.
	Authenticate(token string) (string, error)
}

type UrlService interface {
	CreateUrl(owner uuid.UUID, originalUrl, shortenedUrl string) (entity.URL, error)
	CreateUrlShortCode(originalUrl string) (entity.URL, error)
	GetByShortUrl(shortenedUrl string) (entity.URL, error)
	GetByOwner(owner uuid.UUID) ([]entity.URL, error)
	GetByKeyword(keyword string) ([]entity.URL, error)
	GetByKeywords(keywords []string) ([]entity.URL, error)
	GetByOriginalUrl(originalUrl string) ([]entity.URL, error)
	GetById(id uuid.UUID) (entity.URL, error)
	Remove(id uuid.UUID) error
}

type UserService interface {
	CreateUser(email, password string) (entity.User, error)
	GetUserByEmail(email string) (entity.User, error)
	GetUserByID(id uuid.UUID) (uuid.UUID, error)
	GetUserByToken(token string) (uuid.UUID, error)
	GetUserByUsername(username string) (uuid.UUID, error)
	RemoveUser(id uuid.UUID) error
}
