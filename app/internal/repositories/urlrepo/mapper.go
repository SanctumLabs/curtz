package urlRepo

import (
	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/internal/repositories/models"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

func mapModelToEntity(url models.Url) entities.URL {
	keywords := []entities.Keyword{}

	for _, kw := range url.Keywords {
		keywords = append(keywords, entities.Keyword{
			ID:    identifier.New().FromString(kw.UrlId),
			Value: kw.Value,
		})
	}

	return entities.URL{
		ID:          identifier.New().FromString(url.BaseModel.Id),
		UserId:      identifier.New().FromString(url.UserId),
		ShortCode:   url.ShortCode,
		CustomAlias: url.CustomAlias,
		OriginalUrl: url.OriginalURL,
		Hits:        uint(url.VisitCount),
		ExpiresOn:   url.ExpiresOn,
		Keywords:    keywords,
		BaseEntity: entities.BaseEntity{
			CreatedAt: url.BaseModel.CreatedAt,
			DeletedAt: url.BaseModel.DeletedAt,
			UpdatedAt: url.BaseModel.UpdatedAt,
		},
	}
}

func mapEntityToModel(url entities.URL) models.Url {
	keywords := make([]models.Keyword, len(url.Keywords))

	for i, keyword := range url.Keywords {
		keywords[i] = models.Keyword{
			UrlId: url.ID.String(),
			Value: keyword.Value,
		}
	}

	return models.Url{
		BaseModel: models.BaseModel{
			Id:        url.ID.String(),
			CreatedAt: url.CreatedAt,
			UpdatedAt: url.UpdatedAt,
		},
		UserId:      url.UserId.String(),
		OriginalURL: url.OriginalUrl,
		ShortCode:   url.ShortCode,
		CustomAlias: url.CustomAlias,
		ExpiresOn:   url.ExpiresOn,
		Keywords:    keywords,
	}
}
