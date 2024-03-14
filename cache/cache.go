package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Cache struct {
	client *redis.Client
}

func NewCache(addr string) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr: addr, // e.g., "localhost:6379"
	})

	return &Cache{client: client}
}

func (c *Cache) Get(key string) (string, bool) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", false
	}
	return val, true
}

func (c *Cache) Set(key string, value string, ttl int) {
	c.client.Set(ctx, key, value, time.Duration(ttl)*time.Second)
}
