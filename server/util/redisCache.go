package util

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	Client *redis.Client
	Ctx    context.Context
	Key    string
	Ttl    time.Duration
}

const DEFAULT_TTL = "10"

func NewRedisCache(client *redis.Client, key string) RedisCache {
	ttl, _ := strconv.Atoi(GetEnv("CACHE_TTL", DEFAULT_TTL))
	return RedisCache{
		client, context.Background(), key, time.Duration(ttl) * time.Minute,
	}
}
