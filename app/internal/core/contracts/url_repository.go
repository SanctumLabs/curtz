package contracts

import "github.com/sanctumlabs/curtz/app/internal/core/entities"

type UrlRepository interface {
	UrlReadRepository
	UrlWriteRepository
}

type UrlWriteRepository interface {
	Save(entities.URL) (entities.URL, error)
	Delete(id string) error
	IncrementHits(shortCode string) error
}

type UrlReadRepository interface {
	GetByShortCode(shortCode string) (entities.URL, error)
	GetByOwner(owner string) ([]entities.URL, error)
	GetByKeyword(keyword string) ([]entities.URL, error)
	GetByKeywords(keywords []string) ([]entities.URL, error)
	GetByOriginalUrl(originalUrl string) (entities.URL, error)
	GetById(id string) (entities.URL, error)
}
