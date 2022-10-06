package data

import (
	"time"

	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

// MockUser creates a mock user given an email and a password
func MockUser(email, password string) (entities.User, error) {
	user, err := entities.NewUser(email, password)
	return user, err
}

func MockUrl(userID, originalUrl, customAlias, shortCode string, expiresOn time.Time, keyWords []string) entities.URL {
	mockUrl, err := entities.NewUrl(identifier.New().FromString(userID), originalUrl, customAlias, expiresOn, keyWords)
	if err != nil {
		return entities.URL{}
	}
	err = mockUrl.SetShortCode(shortCode)
	if err != nil {
		return entities.URL{}
	}

	return *mockUrl
}
