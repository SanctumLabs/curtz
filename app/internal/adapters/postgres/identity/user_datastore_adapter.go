package identitydatastore

import (
	"context"

	"github.com/google/wire"
	postgresrepo "github.com/sanctumlabs/curtz/app/internal/adapters/postgres"
	"github.com/sanctumlabs/curtz/app/internal/domain/identity"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
)

type (
	userWriteDatastoreAdapter struct {
		logPrefix string
		dbClient  database.PostgresDatabaseClient
		config    database.Config
		// withTx is the transaction executor. In production it wraps postgres.WithTransaction;
		// in tests it can be replaced with a function that calls the mock querier directly,
		// bypassing the real database entirely.
		withTx func(ctx context.Context, fn func(q postgresrepo.UserWriteQuerier) (identity.User, error)) (identity.User, error)
	}

	userReadDatastoreAdapter struct {
		logPrefix string
		dbClient  database.PostgresDatabaseClient
		config    database.Config
		withTx    func(ctx context.Context, fn func(q postgresrepo.UserReadQuerier) (identity.User, error)) (identity.User, error)
	}

	userDatastoreAdapter struct {
		userReadDatastoreAdapter
		userWriteDatastoreAdapter
	}
)

var (
	_ identity.UserDatastore      = (*userDatastoreAdapter)(nil)
	_ identity.UserWriteDatastore = (*userWriteDatastoreAdapter)(nil)
	_ identity.UserReadDatastore  = (*userReadDatastoreAdapter)(nil)

	UserWriteDatastoreAdapter = wire.NewSet(NewUserWriteDatastoreAdapter)
	UserReadDatastoreAdapter  = wire.NewSet(NewUserReadRepoAdapter)
	UserDatastoreAdapter      = wire.NewSet(NewUserDatastoreAdapter)
)

func NewUserDatastoreAdapter(dbClient database.PostgresDatabaseClient, config database.Config) identity.UserDatastore {
	repo := &userDatastoreAdapter{
		userReadDatastoreAdapter:  *NewUserReadRepoAdapter(dbClient, config).(*userReadDatastoreAdapter),
		userWriteDatastoreAdapter: *NewUserWriteDatastoreAdapter(dbClient, config).(*userWriteDatastoreAdapter),
	}

	return repo
}
