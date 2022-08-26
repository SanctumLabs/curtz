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

func MockUrl(userID, originalUrl, customAlias string, expiresOn time.Time, keyWords []string) (*entities.URL, error) {
	userId := identifier.New().FromString(userID)
	url, err := entities.NewUrl(userId, originalUrl, customAlias, expiresOn, keyWords)
	return url, err
}
