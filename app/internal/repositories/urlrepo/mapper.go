package urlRepo

import (
	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/internal/repositories/models"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

func mapModelToEntity(url models.Url) entities.URL {
	keywords := []string{}

	for _, kw := range url.Keywords {
		keywords = append(keywords, kw.Value)
	}

	id := identifier.New().FromString(url.BaseModel.Id)
	userId := identifier.New().FromString(url.UserId)

	urlEntity, err := entities.NewUrl(userId, url.OriginalURL, url.CustomAlias, url.ExpiresOn, keywords)
	if err != nil {
		return entities.URL{}
	}

	urlEntity.ID = id
	urlEntity.BaseEntity.CreatedAt = url.BaseModel.CreatedAt
	urlEntity.BaseEntity.UpdatedAt = url.BaseModel.UpdatedAt
	err = urlEntity.SetHits(url.VisitCount)
	if err != nil {
		return entities.URL{}
	}

	return *urlEntity
}

func mapEntityToModel(url entities.URL) models.Url {
	keywords := make([]models.Keyword, len(url.GetKeywords()))

	for i, keyword := range url.GetKeywords() {
		keywords[i] = models.Keyword{
			UrlId: url.ID.String(),
			Value: keyword.GetValue(),
		}
	}

	return models.Url{
		BaseModel: models.BaseModel{
			Id:        url.ID.String(),
			CreatedAt: url.CreatedAt,
			UpdatedAt: url.UpdatedAt,
		},
		UserId:      url.UserId.String(),
		OriginalURL: url.GetOriginalURL(),
		ShortCode:   url.GetShortCode(),
		CustomAlias: url.GetCustomAlias(),
		ExpiresOn:   url.GetExpiresOn(),
		Keywords:    keywords,
	}
}
