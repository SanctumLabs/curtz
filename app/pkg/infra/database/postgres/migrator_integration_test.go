//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sanctumlabs/curtz/app/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Migrations are applied through Migrate, record their state in schema_migrations (the table `make migrate` and the
// compose migrate job use) and a second run changes nothing.
func TestMigrate_AppliesTheMigrationsOnceInSchemaMigrations(t *testing.T) {
	ctx := context.Background()
	container, err := test.TestPostgresDatabaseContainer(ctx, test.DefaultTestDatabaseConfig())
	require.NoError(t, err)
	t.Cleanup(func() { _ = container.Terminate(context.Background()) })
	connectionString, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	require.NoError(t, test.RunDatabaseMigration(ctx, connectionString))
	require.NoError(t, test.RunDatabaseMigration(ctx, connectionString), "a second run must be a no-op")

	pool, err := pgxpool.New(ctx, connectionString)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	var version int64
	var dirty bool
	require.NoError(t, pool.QueryRow(ctx, "SELECT version, dirty FROM schema_migrations").Scan(&version, &dirty))
	assert.GreaterOrEqual(t, version, int64(1))
	assert.False(t, dirty)

	var legacyTable, outboxTable *string
	require.NoError(t, pool.QueryRow(ctx, "SELECT to_regclass('bid_schema_migrations')::text").Scan(&legacyTable))
	require.NoError(t, pool.QueryRow(ctx, "SELECT to_regclass('outbox_events')::text").Scan(&outboxTable))
	assert.Nil(t, legacyTable, "the old bid_schema_migrations table must not be created")
	assert.NotNil(t, outboxTable, "the migrations must have created outbox_events")
}
