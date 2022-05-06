package contracts

import (
	"github.com/sanctumlabs/curtz/app/internal/core/domain/entity"
)

type AuthService interface {
	// Authenticate a user given the token. Returns user id if authenticated, error otherwise.
	Authenticate(token string) (string, error)
}

type UrlService interface {
	CreateUrl(owner string, originalUrl, shortenedUrl string) (entity.URL, error)
	CreateUrlShortCode(originalUrl string) (entity.URL, error)
	GetByShortUrl(shortenedUrl string) (entity.URL, error)
	GetByOwner(owner string) ([]entity.URL, error)
	GetByKeyword(keyword string) ([]entity.URL, error)
	GetByKeywords(keywords []string) ([]entity.URL, error)
	GetByOriginalUrl(originalUrl string) ([]entity.URL, error)
	GetById(id string) (entity.URL, error)
	Remove(id string) error
}

type UserService interface {
	CreateUser(email, password string) (entity.User, error)
	GetUserByEmail(email string) (entity.User, error)
	GetUserByID(id string) (entity.User, error)
	GetUserByToken(token string) (entity.User, error)
	GetUserByUsername(username string) (entity.User, error)
	RemoveUser(id string) error
}
