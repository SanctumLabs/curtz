package urlrepo

import (
	"github.com/google/wire"
	"github.com/sanctumlabs/curtz/app/internal/domain/url"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
)

type (
	urlWriteRepositoryAdapter struct {
		logPrefix string
		dbClient  database.PostgresDatabaseClient
	}

	urlReadRepositoryAdapter struct {
		logPrefix string
		dbClient  database.PostgresDatabaseClient
	}

	urlRepoAdapter struct {
		urlReadRepositoryAdapter
		urlWriteRepositoryAdapter
	}
)

var (
	_ url.UrlRepository      = (*urlRepoAdapter)(nil)
	_ url.UrlWriteRepository = (*urlWriteRepositoryAdapter)(nil)
	_ url.UrlReadRepository  = (*urlReadRepositoryAdapter)(nil)

	UrlWriteRepoAdapter = wire.NewSet(NewUrlWriteRepoAdapter)
	UrlReadRepoAdapter  = wire.NewSet(NewUrlReadRepoAdapter)
	UrlRepoAdapter      = wire.NewSet(NewUrlRepoAdapter)
)

func NewUrlRepoAdapter(dbClient database.PostgresDatabaseClient) url.UrlRepository {
	repo := &urlRepoAdapter{
		urlReadRepositoryAdapter:  *NewUrlReadRepoAdapter(dbClient).(*urlReadRepositoryAdapter),
		urlWriteRepositoryAdapter: *NewUrlWriteRepoAdapter(dbClient).(*urlWriteRepositoryAdapter),
	}

	return repo
}
