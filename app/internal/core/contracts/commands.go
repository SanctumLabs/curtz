package contracts

import (
	"time"

	"github.com/sanctumlabs/curtz/app/pkg/validators"
)

type CreateUrlCommand struct {
	userID      string
	originalUrl string
	customAlias string
	expiresOn   time.Time
	keywords    []string
}

type UpdateUrlCommand struct {
	UserId      string
	UrlId       string
	CustomAlias string
	Keywords    []string
	ExpiresOn   *time.Time
}

func NewUpdateUrlCommand(userId, urlId, customAlias string, keywords []string, expiresOn *time.Time) (UpdateUrlCommand, error) {
	if err := validators.IsValidUserId(userId); err != nil {
		return UpdateUrlCommand{}, err
	}

	if err := validators.IsValidUrlId(urlId); err != nil {
		return UpdateUrlCommand{}, err
	}

	if len(customAlias) != 0 {
		if err := validators.IsValidCustomAlias(customAlias); err != nil {
			return UpdateUrlCommand{}, err
		}
	}

	if expiresOn != nil {
		if err := validators.IsValidExpirationTime(*expiresOn); err != nil {
			return UpdateUrlCommand{}, err
		}
	}

	return UpdateUrlCommand{
		UserId:      userId,
		UrlId:       urlId,
		CustomAlias: customAlias,
		Keywords:    keywords,
		ExpiresOn:   expiresOn,
	}, nil
}
