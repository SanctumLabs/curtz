package urlsvc

import (
	"time"

	"github.com/sanctumlabs/curtz/app/internal/core/contracts"
	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

//UrlSvc represents a url service use case
type UrlWriteSvc struct {
	// repo is an interface used to perform CUD operations on URL records. Does not perform write operations
	repo contracts.UrlWriteRepository
	// userSvc is an interface used to interact with the user service use cases
	userSvc contracts.UserService
}

// NewUrlSvc creates a new url service
func NewUrlWriteSvc(urlRepository contracts.UrlWriteRepository, userSvc contracts.UserService) *UrlWriteSvc {
	return &UrlWriteSvc{urlRepository, userSvc}
}

// CreateUrl creates a new shorted url given a user id, original url, custom alias, when it should expire and slice of keywords
func (svc *UrlWriteSvc) CreateUrl(userId string, originalUrl string, customAlias string, expiresOn time.Time, keywords []string) (entities.URL, error) {
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

// Remove removes a saved shortened url given its ID
func (svc *UrlWriteSvc) Remove(id string) error {
	err := svc.repo.Delete(id)

	if err != nil {
		return err
	}
	return nil
}
