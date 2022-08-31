package urlRepo

import (
	"context"

	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/internal/repositories/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UrlReadRepo struct {
	dbClient *mongo.Collection
	ctx      context.Context
	UrlRepo
}

func NewUrlReadRepo(dbClient *mongo.Collection, ctx context.Context) *UrlReadRepo {
	repo := &UrlReadRepo{
		dbClient: dbClient,
		ctx:      ctx,
		UrlRepo:  *NewUrlRepo(dbClient, ctx),
	}

	return repo
}

func (r *UrlReadRepo) GetByShortCode(shortCode string) (entities.URL, error) {
	url, err := r.getSingleResult("short_code", shortCode)
	if err != nil {
		return entities.URL{}, err
	}

	return url, nil
}

func (r *UrlReadRepo) GetByOwner(owner string) ([]entities.URL, error) {
	urls := []entities.URL{}
	filter := bson.D{{Key: "user_id", Value: owner}}

	cursor, err := r.dbClient.Find(r.ctx, filter)
	if err != nil {
		return nil, err
	}

	var results []bson.D
	if err = cursor.All(r.ctx, &results); err != nil {
		return nil, err
	}

	for _, result := range results {
		document, err := bson.Marshal(result)
		if err != nil {
			return nil, err
		}

		var url models.Url
		err = bson.Unmarshal(document, &url)
		if err != nil {
			return nil, err
		}
		urlEntity := mapModelToEntity(url)
		urls = append(urls, urlEntity)
	}

	return urls, nil
}

func (r *UrlReadRepo) GetByKeyword(keyword string) ([]entities.URL, error) {
	panic("implement me")
}

func (r *UrlReadRepo) GetByKeywords(keywords []string) ([]entities.URL, error) {
	panic("implement me")
}

func (r *UrlReadRepo) GetByOriginalUrl(originalUrl string) (entities.URL, error) {
	url, err := r.getSingleResult("original_url", originalUrl)
	if err != nil {
		return entities.URL{}, err
	}

	return url, nil
}

func (r *UrlReadRepo) GetById(id string) (entities.URL, error) {
	url, err := r.getSingleResult("id", id)

	if err != nil {
		return entities.URL{}, err
	}

	return url, nil
}
