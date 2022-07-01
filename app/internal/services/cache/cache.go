package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/sanctumlabs/curtz/app/config"
)

type Cache struct {
	client redis.UniversalClient
	ctx    context.Context
}

func New(config config.CacheConfig) *Cache {
	ctx := context.Background()
	var options redis.UniversalOptions

	if config.RequireAuth {
		options = redis.UniversalOptions{
			Password: config.Password,
			Username: config.Username,
			Addrs:    []string{fmt.Sprintf(":%s", config.Port)},
		}
	} else {
		options = redis.UniversalOptions{
			Addrs: []string{fmt.Sprintf(":%s", config.Port)},
		}
	}

	redisClient := redis.NewUniversalClient(&options)

	return &Cache{
		client: redisClient,
		ctx:    ctx,
	}
}

func (c *Cache) LookupUrl(shortCode string) (string, error) {
	originalUrl, err := c.client.Get(c.ctx, shortCode).Result()

	if err != nil {
		return "", err
	}

	return originalUrl, nil
}

func (c *Cache) SaveUrl(shortCode, originalUrl string) (string, error) {
	cmd, err := c.client.Set(c.ctx, shortCode, originalUrl, 0).Result()

	if err != nil {
		return "", err
	}

	return cmd, nil
}
