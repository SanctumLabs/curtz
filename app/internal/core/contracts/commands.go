package contracts

import "time"

type CreateUrlCommand struct {
	userID      string
	originalUrl string
	customAlias string
	expiresOn   time.Time
	keywords    []string
}

type UpdateUrlCommand struct {
	userId      string
	urlId       string
	customAlias string
	keywords    []string
	expiresOn   time.Time
}

func NewUpdateUrlCommand(userId, urlId, customAlias string, keywords []string, expiresOn time.Time) (UpdateUrlCommand, error) {

	return UpdateUrlCommand{
		userId:      userId,
		urlId:       urlId,
		customAlias: customAlias,
		keywords:    keywords,
		expiresOn:   expiresOn,
	}, nil
}
