package url

import (
	"github.com/google/uuid"
	"github.com/sanctumlabs/curtz/internal/core/contracts"
	"github.com/sanctumlabs/curtz/internal/core/domain/url"
)

type UseCase struct {
	urlRepo contracts.UrlRepository
}

func NewUseCase(urlRepository contracts.UrlRepository) *UseCase {
	return &UseCase{urlRepository}
}

func (s *UseCase) CreateUrl(owner uuid.UUID, originalUrl, shortenedUrl string) (url.URL, error) {
	return s.urlRepo.CreateUrl(owner, originalUrl, shortenedUrl)
}
