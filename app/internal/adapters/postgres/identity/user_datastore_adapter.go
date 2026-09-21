package identitydatastore

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/wire"
	postgresrepo "github.com/sanctumlabs/curtz/app/internal/adapters/postgres"
	"github.com/sanctumlabs/curtz/app/internal/domain/identity"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
	recoveryutils "github.com/sanctumlabs/curtz/app/pkg/utils/recover"
)

type (
	// txExecutor runs fn with a querier bound to a transaction. In production it wraps
	// postgres.WithTransactionVoid; in tests it can be replaced with a function that calls
	// the mock querier directly, bypassing the real database entirely.
	txExecutor[Q any] func(ctx context.Context, fn func(q Q) error) error

	userWriteDatastoreAdapter struct {
		logPrefix string
		dbClient  database.PostgresDatabaseClient
		config    database.Config
		withTx    txExecutor[postgresrepo.UserWriteQuerier]
	}

	userReadDatastoreAdapter struct {
		logPrefix string
		dbClient  database.PostgresDatabaseClient
		config    database.Config
		withTx    txExecutor[postgresrepo.UserReadQuerier]
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

// execute runs fn inside a transaction, bounded by the configured operation timeout and retried
// according to the retry config. The context is checked before fn runs so a cancelled request
// never reaches the database.
func execute[T any, Q any](
	ctx context.Context,
	config database.Config,
	operationName string,
	withTx txExecutor[Q],
	fn func(ctx context.Context, qtx Q) (T, error),
) (T, error) {
	operationCtx, operationCancel := context.WithTimeout(ctx, config.OperationTimeout)
	defer operationCancel()

	return recoveryutils.ExecuteWithRetry(
		operationCtx,
		func(retryCtx context.Context) (T, error) {
			var result T
			err := withTx(retryCtx, func(qtx Q) error {
				if ctxErr := retryCtx.Err(); ctxErr != nil {
					slog.ErrorContext(retryCtx, "Operation cancelled before validation with error", "error", ctxErr)
					return fmt.Errorf("operation cancelled before validation: %w", ctxErr)
				}
				r, fnErr := fn(retryCtx, qtx)
				if fnErr != nil {
					return fnErr
				}
				result = r
				return nil
			})
			return result, err
		},
		config.RetryConfig,
		operationName,
	)
}
