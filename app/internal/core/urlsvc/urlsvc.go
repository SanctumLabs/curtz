package urlsvc

import (
	"github.com/sanctumlabs/curtz/app/internal/core/contracts"
	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

type UrlSvc struct {
	repo contracts.UrlRepository
}

func NewUrlSvc(urlRepository contracts.UrlRepository) *UrlSvc {
	return &UrlSvc{urlRepository}
}

func (svc *UrlSvc) CreateUrl(userId, originalUrl, customAlias, expiresOn string, keywords []string) (entities.URL, error) {
	userIdentifier := identifier.New().FromString(userId)

	url, err := entities.NewUrl(userIdentifier, originalUrl, customAlias, expiresOn, keywords)

	if err != nil {
		return entities.URL{}, err
	}

	return svc.repo.Save(*url)
}

func (svc *UrlSvc) GetByShortCode(shortCode string) (entities.URL, error) {
	panic("implement me")
}

func (svc *UrlSvc) GetByUserId(userId string) ([]entities.URL, error) {
	panic("implement me")
}

func (svc *UrlSvc) GetByKeyword(keyword string) ([]entities.URL, error) {
	panic("implement me")
}

func (svc *UrlSvc) GetByKeywords(keywords []string) ([]entities.URL, error) {
	panic("implement me")
}

func (svc *UrlSvc) GetByOriginalUrl(originalUrl string) ([]entities.URL, error) {
	panic("implement me")
}

func (svc *UrlSvc) GetById(id string) (entities.URL, error) {
	panic("implement me")
}

func (svc *UrlSvc) Remove(id string) error {
	panic("implement me")
}
