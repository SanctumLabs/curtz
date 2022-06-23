package contracts

import (
	"github.com/sanctumlabs/curtz/app/internal/core/entities"
)

type UrlRepository interface {
	Save(entities.URL) (entities.URL, error)
	GetByShortUrl(shortenedUrl string) (entities.URL, error)
	GetByOwner(owner string) ([]entities.URL, error)
	GetByKeyword(keyword string) ([]entities.URL, error)
	GetByKeywords(keywords []string) ([]entities.URL, error)
	GetByOriginalUrl(originalUrl string) ([]entities.URL, error)
	GetById(id string) (entities.URL, error)
	Delete(id string) error
}

type UserRepository interface {
	CreateUser(entities.User) (entities.User, error)
	GetByEmail(email string) (entities.User, error)
	GetById(id string) (entities.User, error)
	RemoveUser(id string) error
}
