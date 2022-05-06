package contracts

import (
	"github.com/sanctumlabs/curtz/app/internal/core/domain/entity"
)

type UrlRepository interface {
	Save(owner string, originalUrl, shortenedUrl string) (entity.URL, error)
	GetByShortUrl(shortenedUrl string) (entity.URL, error)
	GetByOwner(owner string) ([]entity.URL, error)
	GetByKeyword(keyword string) ([]entity.URL, error)
	GetByKeywords(keywords []string) ([]entity.URL, error)
	GetByOriginalUrl(originalUrl string) ([]entity.URL, error)
	GetById(id string) (entity.URL, error)
	Delete(id string) error
}

type UserRepository interface {
	CreateUser(email, password string) (entity.User, error)
	GetByEmail(email string) (entity.User, error)
	GetById(id string) (entity.User, error)
	GetByUsername(username string) (entity.User, error)
	RemoveUser(id string) error
}
