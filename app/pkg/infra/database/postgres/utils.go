package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func buildConnectionString(config PostgresDatabaseConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?pool_max_conns=%d&pool_min_conns=%d",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.MaxConns,
		config.MinConns,
	)
}

// WithTransactionRetry wraps a transaction with retry logic
func WithTransactionRetry[T any](
	ctx context.Context,
	db *pgxpool.Pool,
	operation func(tx pgx.Tx) (T, error),
	config tools.RetryConfig,
	operationName string,
) (T, error) {
	return tools.ExecuteWithRetry(ctx, func(ctx context.Context) (T, error) {
		tx, err := db.Begin(ctx)
		if err != nil {
			var result T
			return result, fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback(ctx) // This is safe even if commit succeeds

		result, err := operation(tx)
		if err != nil {
			return result, err
		}

		if err := tx.Commit(ctx); err != nil {
			return result, fmt.Errorf("failed to commit transaction: %w", err)
		}

		return result, nil
	}, config, operationName)
}
