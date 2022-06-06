package urlsvc

import (
	"github.com/sanctumlabs/curtz/app/internal/core/domain"
	entity "github.com/sanctumlabs/curtz/app/internal/core/entities"
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

func (svc *service) Remove(id string) error {
	panic("implement me")
}

func (svc *service) CreateUrl(owner string, originalUrl, shortenedUrl string) (entity.URL, error) {
	return svc.useCase.CreateUrl(owner, originalUrl, shortenedUrl)
}

func (svc *service) GetByShortUrl(shortenedUrl string) (entity.URL, error) {
	panic("implement me")
}

func (svc *service) GetByOwner(owner string) ([]entity.URL, error) {
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

func (svc *service) GetById(id string) (entity.URL, error) {
	panic("implement me")
}
