package url

import (
	"github.com/google/uuid"
	"github.com/sanctumlabs/curtz/internal/core/contracts"
	"github.com/sanctumlabs/curtz/internal/core/domain/url"
)

type UrlUseCase struct {
	urlRepo contracts.UrlRepository
}

func NewUrlUseCase(urlRepository contracts.UrlRepository) *UrlUseCase {
	return &UrlUseCase{urlRepository}
}

func (s *UrlUseCase) CreateUrl(owner uuid.UUID, originalUrl, shortenedUrl string) (url.URL, error) {
	return s.urlRepo.CreateUrl(owner, originalUrl, shortenedUrl)
}
