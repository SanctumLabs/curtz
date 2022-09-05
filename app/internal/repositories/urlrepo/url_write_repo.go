package urlRepo

import (
	"context"
	"time"

	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/internal/repositories/models"
	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UrlWriteRepo struct {
	dbClient *mongo.Collection
	ctx      context.Context
	UrlRepo
}

func NewUrlWriteRepo(dbClient *mongo.Collection, ctx context.Context) *UrlWriteRepo {
	repo := &UrlWriteRepo{
		dbClient: dbClient,
		ctx:      ctx,
		UrlRepo:  *NewUrlRepo(dbClient, ctx),
	}

	return repo
}

func (r *UrlWriteRepo) Save(url entities.URL) (entities.URL, error) {
	if _, err := r.getSingleResult("original_url", url.OriginalUrl); err == nil {
		return entities.URL{}, errdefs.ErrURLAlreadyExists
	}

	urlModel := mapEntityToModel(url)

	_, err := r.dbClient.InsertOne(r.ctx, urlModel)

	if err != nil {
		return url, err
	}

	return url, nil
}

// Update performs an update on an existing shortened URL given urlID, customAlias, keywords & expiresOn
func (r *UrlWriteRepo) Update(urlID, customAlias string, keywords []entities.Keyword, expiresOn time.Time) (entities.URL, error) {
	existingUrl, err := r.getSingleResult("id", urlID)
	if err != nil {
		return entities.URL{}, errdefs.ErrURLAlreadyExists
	}

	existingUrl.CustomAlias = customAlias
	existingUrl.Keywords = append(existingUrl.Keywords, entities.Keyword{})
	existingUrl.ExpiresOn = expiresOn

	kws := make([]models.Keyword, len(keywords))
	if len(keywords) != 0 {
		for i, keyword := range keywords {
			kws[i] = models.Keyword{
				UrlId: urlID,
				Value: keyword.Value,
			}
		}
	}

	filter := bson.D{{Key: "id", Value: urlID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "custom_alias", Value: customAlias}}},
		{Key: "$set", Value: bson.D{{Key: "expires_on", Value: expiresOn}}},
		{Key: "$addToSet", Value: bson.D{{Key: "keywords", Value: kws}}},
	}

	opts := options.Update().SetUpsert(false)

	_, err = r.dbClient.UpdateOne(r.ctx, filter, update, opts)
	if err != nil {
		return entities.URL{}, err
	}

	return existingUrl, nil
}

// Delete deletes a url given its ID
func (r *UrlWriteRepo) Delete(id string) error {
	if _, err := r.getSingleResult("id", id); err != nil {
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

func (r *UrlWriteRepo) IncrementHits(shortCode string) error {
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
