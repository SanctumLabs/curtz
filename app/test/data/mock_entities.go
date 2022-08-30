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
	return entities.URL{
		ID:          identifier.New(),
		UserId:      identifier.New().FromString(userID),
		OriginalUrl: originalUrl,
		CustomAlias: customAlias,
		ExpiresOn:   expiresOn,
		BaseEntity: entities.BaseEntity{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Hits:      0,
		Keywords:  []entities.Keyword{},
		ShortCode: shortCode,
	}
}
