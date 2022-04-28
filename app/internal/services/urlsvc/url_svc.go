package urlsvc

import (
	"github.com/google/uuid"
	"github.com/sanctumlabs/curtz/app/internal/core/domain"
	"github.com/sanctumlabs/curtz/app/internal/core/domain/entity"
)

type service struct {
	useCase *domain.UrlInteractor
}

func NewUrlService(urlInteractor *domain.UrlInteractor) *service {
	return &service{useCase: urlInteractor}
}

func (svc *service) CreateUrlShortCode(originalUrl string) (entity.URL, error) {
	return svc.useCase.CreateShortCode(originalUrl)
}

func (svc *service) Remove(id uuid.UUID) error {
	panic("implement me")
}

func (svc *service) CreateUrl(owner uuid.UUID, originalUrl, shortenedUrl string) (entity.URL, error) {
	return svc.useCase.CreateUrl(owner, originalUrl, shortenedUrl)
}

func (svc *service) GetByShortUrl(shortenedUrl string) (entity.URL, error) {
	panic("implement me")
}

func (svc *service) GetByOwner(owner uuid.UUID) ([]entity.URL, error) {
	panic("implement me")
}

func (svc *service) GetByKeyword(keyword string) ([]entity.URL, error) {
	panic("implement me")
}

func (svc *service) GetByKeywords(keywords []string) ([]entity.URL, error) {
	panic("implement me")
}

func (svc *service) GetByOriginalUrl(originalUrl string) ([]entity.URL, error) {
	panic("implement me")
}

func (svc *service) GetById(id uuid.UUID) (entity.URL, error) {
	panic("implement me")
}
