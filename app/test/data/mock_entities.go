package data

import (
	"time"

	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

// MockUser creates a mock user given an email and a password
func MockUser(email, password string) (entities.User, error) {
	user, err := entities.NewUser(email, password)
	return *user, err
}

func MockUrl(userID, originalUrl, customAlias, shortCode string, expiresOn time.Time, keyWords []string) *entities.URL {
	id, err := identifier.New().FromString(userID)
	if err != nil {
		return nil
	}

	mockUrl, err := entities.NewUrl(id, originalUrl, customAlias, expiresOn, keyWords)
	if err != nil {
		return nil
	}

	err = mockUrl.SetShortCode(shortCode)
	if err != nil {
		return nil
	}

	return mockUrl
}
