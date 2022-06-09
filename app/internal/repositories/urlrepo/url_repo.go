package urlRepo

import (
	"context"

	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/internal/repositories/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UrlRepo struct {
	dbClient *mongo.Collection
	ctx      context.Context
}

func NewUrlRepo(dbClient *mongo.Collection, ctx context.Context) *UrlRepo {
	repo := &UrlRepo{
		dbClient: dbClient,
		ctx:      ctx,
	}

	return repo
}

func (r *UrlRepo) Save(url entities.URL) (entities.URL, error) {
	keywords := make([]models.Keyword, len(url.Keywords))

	for i, keyword := range url.Keywords {
		keywords[i] = models.Keyword{
			UrlId: url.ID.String(),
			Value: keyword.Value,
		}
	}

	urlModel := models.Url{
		BaseModel: models.BaseModel{
			Id:        url.ID.String(),
			CreatedAt: url.CreatedAt,
		},
		UserId:      url.UserId.String(),
		OriginalURL: url.OriginalUrl,
		ShortCode:   url.ShortCode,
		CustomAlias: url.CustomAlias,
		ExpiresOn:   url.ExpiresOn,
		Keywords:    keywords,
	}

	_, err := r.dbClient.InsertOne(context.TODO(), urlModel)

	if err != nil {
		return url, err
	}

	return url, nil
}

func (r *UrlRepo) GetByShortUrl(shortenedUrl string) (entities.URL, error) {
	panic("implement me")
}

func (r *UrlRepo) GetByOwner(owner string) ([]entities.URL, error) {
	panic("implement me")
}

func (r *UrlRepo) GetByKeyword(keyword string) ([]entities.URL, error) {
	panic("implement me")
}

func (r *UrlRepo) GetByKeywords(keywords []string) ([]entities.URL, error) {
	panic("implement me")
}

func (r *UrlRepo) GetByOriginalUrl(originalUrl string) ([]entities.URL, error) {
	panic("implement me")
}

func (r *UrlRepo) GetById(id string) (entities.URL, error) {
	panic("implement me")
}

func (r *UrlRepo) Delete(id string) error {
	panic("implement me")
}
