package urlrepo

import (
	"context"

	"github.com/sanctumlabs/curtz/app/internal/domain/url"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
)

func NewUrlWriteRepoAdapter(dbClient database.PostgresDatabaseClient) url.UrlWriteRepository {
	repo := &urlWriteRepositoryAdapter{
		dbClient:  dbClient,
		logPrefix: "UrlWriteRepoAdapter",
	}

	return repo
}

func (repo *urlWriteRepositoryAdapter) Create(ctx context.Context, urlEntity url.URL) (url.URL, error) {
	return urlEntity, nil
}

func (repo *urlWriteRepositoryAdapter) Update(ctx context.Context, urlEntity url.URL) (url.URL, error) {
	return urlEntity, nil
}

func (repo *urlWriteRepositoryAdapter) SoftDelete(ctx context.Context, id string) error {
	panic("not implemented")
}

// Delete deletes a given entity by its ID
func (repo *urlWriteRepositoryAdapter) Delete(ctx context.Context, id string) error {
	panic("not implemented")
}
