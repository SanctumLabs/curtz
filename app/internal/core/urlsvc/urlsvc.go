package urlsvc

import (
	"time"

	"github.com/sanctumlabs/curtz/app/internal/core/contracts"
	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

//UrlSvc represents a url service use case
type UrlSvc struct {
	// repo is an interface used to perform CRUD operations on URL records
	repo contracts.UrlRepository
	// userSvc is an interface used to interact with the user service use cases
	userSvc contracts.UserService
	// cache is an interface used to interact with cache service
	cache contracts.CacheService
}

// NewUrlSvc creates a new url service
func NewUrlSvc(urlRepository contracts.UrlRepository, userSvc contracts.UserService, cacheSvc contracts.CacheService) *UrlSvc {
	return &UrlSvc{urlRepository, userSvc, cacheSvc}
}

// CreateUrl creates a new shorted url given a user id, original url, custom alias, when it should expire and slice of keywords
func (svc *UrlSvc) CreateUrl(userId string, originalUrl string, customAlias string, expiresOn time.Time, keywords []string) (entities.URL, error) {
	if _, err := svc.userSvc.GetUserByID(userId); err != nil {
		return entities.URL{}, err
	}

	userIdentifier := identifier.New().FromString(userId)
	url, err := entities.NewUrl(userIdentifier, originalUrl, customAlias, expiresOn, keywords)

	if err != nil {
		return entities.URL{}, err
	}

	return svc.repo.Save(*url)
}

// GetByShortCode returns shortened url given its short code
func (svc *UrlSvc) GetByShortCode(shortCode string) (entities.URL, error) {
	url, err := svc.repo.GetByShortCode(shortCode)

	if err != nil {
		return entities.URL{}, err
	}

	return url, nil
}

// GetByUserId retrieves all urls for a given user
func (svc *UrlSvc) GetByUserId(userId string) ([]entities.URL, error) {
	if _, err := svc.userSvc.GetUserByID(userId); err != nil {
		return nil, err
	}

	urls, err := svc.repo.GetByOwner(userId)

	if err != nil {
		return nil, err
	}

	return urls, nil
}

// GetByKeyword retrieves url records given a keyword
func (svc *UrlSvc) GetByKeyword(keyword string) ([]entities.URL, error) {
	urls, err := svc.repo.GetByKeyword(keyword)

	if err != nil {
		return nil, err
	}

	return urls, nil
}

// GetByKeywords retrieves url records given their keywords
func (svc *UrlSvc) GetByKeywords(keywords []string) ([]entities.URL, error) {
	urls, err := svc.repo.GetByKeywords(keywords)

	if err != nil {
		return nil, err
	}

	return urls, nil
}

// GetByOriginalUrl retrieves a url given its original url
func (svc *UrlSvc) GetByOriginalUrl(originalUrl string) (entities.URL, error) {
	url, err := svc.repo.GetByOriginalUrl(originalUrl)

	if err != nil {
		return entities.URL{}, err
	}

	return url, nil
}

// GetById retrieves url given its id
func (svc *UrlSvc) GetById(id string) (entities.URL, error) {
	url, err := svc.repo.GetById(id)

	if err != nil {
		return entities.URL{}, err
	}

	return url, nil
}

// Remove removes a saved shortened url given its ID
func (svc *UrlSvc) Remove(id string) error {
	err := svc.repo.Delete(id)

	if err != nil {
		return err
	}
	return nil
}

// LookupUrl looks up the original url given the short code
func (svc *UrlSvc) LookupUrl(shortCode string) (string, error) {
	cachedOriginalUrl, err := svc.cache.LookupUrl(shortCode)

	if err != nil || cachedOriginalUrl == "" {
		url, err := svc.repo.GetByShortCode(shortCode)
		if err != nil {
			return "", err
		}
		// nolint
		go svc.cache.SaveUrl(shortCode, url.OriginalUrl)

		// nolint
		go svc.repo.IncrementHits(shortCode)

		return url.OriginalUrl, nil
	}

	// nolint
	go svc.repo.IncrementHits(shortCode)
	return cachedOriginalUrl, nil
}
