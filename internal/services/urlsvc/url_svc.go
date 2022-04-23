package urlsvc

import (
	"github.com/google/uuid"
	"github.com/sanctumlabs/curtz/internal/core/contracts"
	"github.com/sanctumlabs/curtz/internal/core/domain/url"
)

type service struct {
	repo contracts.UrlRepository
}

func NewUrlService(urlRepository contracts.UrlRepository) *service {
	return &service{repo: urlRepository}
}

func (svc *service) Remove(id uuid.UUID) error {
	panic("implement me")
}

func (svc *service) CreateUrl(owner uuid.UUID, originalUrl, shortenedUrl string) (url.URL, error) {
	return svc.repo.CreateUrl(owner, originalUrl, shortenedUrl)
}

func (svc *service) GetByShortUrl(shortenedUrl string) (url.URL, error) {
	panic("implement me")
}

func (svc *service) GetByOwner(owner uuid.UUID) ([]url.URL, error) {
	panic("implement me")
}

func (svc *service) GetByKeyword(keyword string) ([]url.URL, error) {
	panic("implement me")
}

func (svc *service) GetByKeywords(keywords []string) ([]url.URL, error) {
	panic("implement me")
}

func (svc *service) GetByOriginalUrl(originalUrl string) ([]url.URL, error) {
	panic("implement me")
}

func (svc *service) GetById(id uuid.UUID) (url.URL, error) {
	panic("implement me")
}
