package urlsvc

import (
	"github.com/sanctumlabs/curtz/app/internal/core/contracts"
)

//UrlSvc represents a url service use case
type UrlSvc struct {
	// urlReadRepo is an interface used to perform R operations on URL records
	urlReadRepo contracts.UrlReadRepository
	// urlWriteRepo is an interface that performs CRD operations on URL records
	urlWriteRepo contracts.UrlWriteRepository
	// userSvc is an interface used to interact with the user service use cases
	userSvc contracts.UserService
	// cache is an interface used to interact with cache service
	cache contracts.CacheService
}

// NewUrlSvc creates a new url service
func NewUrlSvc(urlReadRepo contracts.UrlReadRepository, urlWriteRepo contracts.UrlWriteRepository, userSvc contracts.UserService, cacheSvc contracts.CacheService) *UrlSvc {
	return &UrlSvc{urlReadRepo, urlWriteRepo, userSvc, cacheSvc}
}

// LookupUrl looks up the original url given the short code
func (svc *UrlSvc) LookupUrl(shortCode string) (string, error) {
	cachedOriginalUrl, err := svc.cache.LookupUrl(shortCode)

	if err != nil || cachedOriginalUrl == "" {
		url, err := svc.urlReadRepo.GetByShortCode(shortCode)
		if err != nil {
			return "", err
		}
		// nolint
		go svc.cache.SaveUrl(shortCode, url.OriginalUrl)

		// nolint
		go svc.urlWriteRepo.IncrementHits(shortCode)

		return url.OriginalUrl, nil
	}

	// nolint
	go svc.urlWriteRepo.IncrementHits(shortCode)
	return cachedOriginalUrl, nil
}
