package main

import (
	"context"
	"testing"
	"time"

	"github.com/sanctumlabs/curtz/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func lookupOf(env map[string]string) config.Lookup {
	return func(key string) (string, bool) {
		value, ok := env[key]
		return value, ok
	}
}

// Postgres is the one required dependency: with it unreachable the process must exit with an error, not hang or serve.
func TestRun_ReturnsAnErrorWhenPostgresIsUnreachable(t *testing.T) {
	if testing.Short() {
		t.Skip("takes about three seconds of connection retries")
	}
	cfg, err := config.Load(lookupOf(map[string]string{
		"DATABASE_PORT": "1", "DATABASE_CONN_TIMEOUT": "1", "REDIS_ADDRESS": "127.0.0.1:1",
	}))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = run(ctx, cfg)

	require.Error(t, err)
	assert.ErrorContains(t, err, "postgres")
}
