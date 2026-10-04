//go:build integration

package redis_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/sanctumlabs/curtz/app/pkg/infra/cache"
	cacheredis "github.com/sanctumlabs/curtz/app/pkg/infra/cache/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
)

func startRedis(t *testing.T) (cache.CacheClient, *redisContainer.RedisContainer) {
	t.Helper()
	ctx := context.Background()

	container, err := redisContainer.Run(ctx, "redis:7-alpine")
	require.NoError(t, err)
	t.Cleanup(func() { _ = container.Terminate(context.Background()) })

	uri, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	client, err := cacheredis.NewRedisClient(cacheredis.RedisClientConfig{Address: []string{strings.TrimPrefix(uri, "redis://")}})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	return client, container
}

func TestRedisClient_SetGetExistsDelete(t *testing.T) {
	ctx := context.Background()
	client, _ := startRedis(t)

	require.NoError(t, client.Ping(ctx))
	require.NoError(t, client.Set(ctx, cache.CacheItem{Key: "short:abc", Value: "https://example.com"}))

	assert.True(t, client.Exists(ctx, "short:abc"))
	item, err := client.Get(ctx, "short:abc")
	require.NoError(t, err)
	assert.Equal(t, "https://example.com", item.Value)

	require.NoError(t, client.Delete(ctx, "short:abc"))
	assert.False(t, client.Exists(ctx, "short:abc"))
	_, err = client.Get(ctx, "short:abc")
	assert.Error(t, err, "a missing key is an error, not an empty item")
}

func TestRedisClient_ItemsExpireAfterTheirTTL(t *testing.T) {
	ctx := context.Background()
	client, _ := startRedis(t)

	require.NoError(t, client.Set(ctx, cache.CacheItem{Key: "k", Value: "v"}, cache.WithTTL(time.Second)))
	require.True(t, client.Exists(ctx, "k"))

	require.Eventually(t, func() bool { return !client.Exists(ctx, "k") }, 5*time.Second, 100*time.Millisecond)
}

func TestRedisClient_PingReportsAStoppedRedis(t *testing.T) {
	ctx := context.Background()
	client, container := startRedis(t)
	require.NoError(t, client.Ping(ctx))

	require.NoError(t, container.Stop(ctx, nil))

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	assert.Error(t, client.Ping(pingCtx))
}
