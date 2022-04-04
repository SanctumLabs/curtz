package urlsvc

import (
	"github.com/google/uuid"
	"github.com/sanctumlabs/curtz/internal/core/domain/url"
	uc "github.com/sanctumlabs/curtz/internal/core/usecases/url"
)

type UrlService struct {
	urlUseCase uc.UseCase
}

func NewUrlService(urlUseCase uc.UseCase) *UrlService {
	return &UrlService{urlUseCase}
}

func (svc *UrlService) CreateUrl(owner uuid.UUID, originalUrl, shortenedUrl string) (url.URL, error) {
	return svc.urlUseCase.CreateUrl(owner, originalUrl, shortenedUrl)
}

func (svc *UrlService) GetByShortUrl(shortenedUrl string) (url.URL, error) {
	panic("implement me")
}

func (svc *UrlService) GetByOwner(owner uuid.UUID) ([]url.URL, error) {
	panic("implement me")
}

func (svc *UrlService) GetByKeyword(keyword string) ([]url.URL, error) {
	panic("implement me")
}

func (svc *UrlService) GetByKeywords(keywords []string) ([]url.URL, error) {
	panic("implement me")
}

func (svc *UrlService) GetByOriginalUrl(originalUrl string) ([]url.URL, error) {
	panic("implement me")
}

func (svc *UrlService) GetById(id uuid.UUID) (url.URL, error) {
	panic("implement me")
}
