package urlRepo

import (
	"context"

	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/internal/repositories/models"
	"github.com/sanctumlabs/curtz/app/tools/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
