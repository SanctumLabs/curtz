package urlRepo

import (
	"context"

	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/internal/repositories/models"
	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"github.com/sanctumlabs/curtz/app/tools/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var log = logger.NewLogger("urlRepo")

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
	if _, err := r.GetByOriginalUrl(url.OriginalUrl); err == nil {
		return entities.URL{}, errdefs.ErrURLAlreadyExists
	}

	urlModel := mapEntityToModel(url)

	_, err := r.dbClient.InsertOne(r.ctx, urlModel)

	if err != nil {
		return url, err
	}

	return url, nil
}

func (r *UrlRepo) GetByShortCode(shortCode string) (entities.URL, error) {
	url, err := r.getSingleResult("short_code", shortCode)
	if err != nil {
		return entities.URL{}, err
	}

	return url, nil
}

func (r *UrlRepo) GetByOwner(owner string) ([]entities.URL, error) {
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

func (r *UrlRepo) GetByKeyword(keyword string) ([]entities.URL, error) {
	panic("implement me")
}

func (r *UrlRepo) GetByKeywords(keywords []string) ([]entities.URL, error) {
	panic("implement me")
}

func (r *UrlRepo) GetByOriginalUrl(originalUrl string) (entities.URL, error) {
	url, err := r.getSingleResult("original_url", originalUrl)
	if err != nil {
		return entities.URL{}, err
	}

	return url, nil
}

func (r *UrlRepo) GetById(id string) (entities.URL, error) {
	url, err := r.getSingleResult("id", id)

	if err != nil {
		return entities.URL{}, err
	}

	return url, nil
}

func (r *UrlRepo) Delete(id string) error {
	_, err := r.getSingleResult("id", id)

	if err != nil {
		return errdefs.ErrURLNotFound
	}

	filter := bson.D{{Key: "id", Value: id}}

	var result bson.D
	if err := r.dbClient.FindOneAndDelete(r.ctx, filter).Decode(&result); err != nil {
		return err
	}

	document, err := bson.Marshal(result)
	if err != nil {
		return err
	}

	var url models.Url
	err = bson.Unmarshal(document, &url)
	if err != nil {
		return err
	}

	return nil
}

func (r *UrlRepo) IncrementHits(shortCode string) error {
	if url, err := r.getSingleResult("short_code", shortCode); err == nil {
		filter := bson.D{{Key: "short_code", Value: shortCode}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "visit_count", Value: url.Hits + 1}}}}
		opts := options.Update().SetUpsert(false)

		result, err := r.dbClient.UpdateOne(r.ctx, filter, update, opts)
		if err != nil {
			return err
		}

		log.Debugf("UrlShortCode %s incremented by 1. Result: %v", shortCode, result.ModifiedCount)
		return nil
	} else {
		return err
	}
}

func (r *UrlRepo) getSingleResult(key, value string) (entities.URL, error) {
	filter := bson.D{{Key: key, Value: value}}

	var result bson.D
	if err := r.dbClient.FindOne(r.ctx, filter).Decode(&result); err != nil {
		return entities.URL{}, err
	}

	document, err := bson.Marshal(result)
	if err != nil {
		return entities.URL{}, err
	}

	var url models.Url
	err = bson.Unmarshal(document, &url)
	if err != nil {
		return entities.URL{}, err
	}

	return mapModelToEntity(url), nil
}
