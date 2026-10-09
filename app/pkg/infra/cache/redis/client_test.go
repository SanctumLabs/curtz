package redis

import (
	"context"
	"testing"
	"time"

	redisGo "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRedisClient_RejectsAnEmptyAddressList(t *testing.T) {
	_, err := NewRedisClient(RedisClientConfig{})

	require.Error(t, err)
}

// The client dials lazily and go-redis reconnects by itself, so building it must succeed while Redis is down.
func TestNewRedisClient_BuildsWithoutIOAndPingReportsAnUnreachableRedis(t *testing.T) {
	client, err := NewRedisClient(RedisClientConfig{Address: []string{"127.0.0.1:1"}})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	assert.Error(t, client.Ping(ctx))
	assert.NoError(t, client.Close())
}

func TestNewRedisClient_PicksTheClientTypeFromTheAddressCount(t *testing.T) {
	single, err := NewRedisClient(RedisClientConfig{Address: []string{"localhost:7001"}})
	require.NoError(t, err)
	assert.IsType(t, &redisGo.Client{}, single.(*redisClient).client)

	cluster, err := NewRedisClient(RedisClientConfig{Address: []string{"localhost:7001", "localhost:7002"}})
	require.NoError(t, err)
	assert.IsType(t, &redisGo.ClusterClient{}, cluster.(*redisClient).client)
}
