package domain

import (
	"github.com/google/uuid"
	"github.com/sanctumlabs/curtz/app/internal/core/contracts"
	"github.com/sanctumlabs/curtz/app/internal/core/domain/entity"
)

type Interactor struct {
	repo contracts.UrlRepository
}

func NewInteractor(urlRepository contracts.UrlRepository) *Interactor {
	return &Interactor{urlRepository}
}

func (useCase Interactor) CreateShortCode(longUrl string) (entity.URL, error) {
	panic("implement me")
}

func (useCase *Interactor) CreateUrl(owner uuid.UUID, originalUrl, shortenedUrl string) (entity.URL, error) {
	return useCase.repo.Save(owner, originalUrl, shortenedUrl)
}
