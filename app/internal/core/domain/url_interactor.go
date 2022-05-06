package domain

import (
	"github.com/sanctumlabs/curtz/app/internal/core/contracts"
	"github.com/sanctumlabs/curtz/app/internal/core/domain/entity"
)

type UrlInteractor struct {
	repo contracts.UrlRepository
}

func NewUrlInteractor(urlRepository contracts.UrlRepository) *UrlInteractor {
	return &UrlInteractor{urlRepository}
}

func (useCase UrlInteractor) CreateShortCode(longUrl string) (entity.URL, error) {
	panic("implement me")
}

func (useCase *UrlInteractor) CreateUrl(owner string, originalUrl, shortenedUrl string) (entity.URL, error) {
	return useCase.repo.Save(owner, originalUrl, shortenedUrl)
}
