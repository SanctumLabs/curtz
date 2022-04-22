package urlRepo

import (
	"github.com/google/uuid"
	"github.com/sanctumlabs/curtz/internal/core/domain/url"
	"gorm.io/gorm"
)

type UrlRepo struct {
	db *gorm.DB
}

func NewUrlRepo(db *gorm.DB) *UrlRepo {
	return &UrlRepo{
		db: db,
	}
}

func (u UrlRepo) CreateUrl(owner uuid.UUID, originalUrl, shortenedUrl string) (url.URL, error) {
	panic("implement me")
}

func (u UrlRepo) GetByShortUrl(shortenedUrl string) (url.URL, error) {

	panic("implement me")
}

func (u UrlRepo) GetByOwner(owner uuid.UUID) ([]url.URL, error) {

	panic("implement me")
}

func (u UrlRepo) GetByKeyword(keyword string) ([]url.URL, error) {

	panic("implement me")
}

func (u UrlRepo) GetByKeywords(keywords []string) ([]url.URL, error) {

	panic("implement me")
}

func (u UrlRepo) GetByOriginalUrl(originalUrl string) ([]url.URL, error) {

	panic("implement me")
}

func (u UrlRepo) GetById(id uuid.UUID) (url.URL, error) {

	panic("implement me")
}

func (u UrlRepo) Delete(id uuid.UUID) error {

	panic("implement me")
}
