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

func (useCase *UrlSvc) CreateUrl(userId, originalUrl, customAlias, expiresOn string, keywords []string) (entities.URL, error) {
	userIdentifier := identifier.New().FromString(userId)

	url, err := entities.NewUrl(userIdentifier, originalUrl, customAlias, expiresOn, keywords)

	if err != nil {
		return entities.URL{}, err
	}

	return useCase.repo.Save(*url)
}

func (useCase UrlSvc) CreateShortCode(longUrl string) (entities.URL, error) {
	panic("implement me")
}
