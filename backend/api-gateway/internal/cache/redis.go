package cache

import (
	"api-gateway/internal/config"
	"api-gateway/internal/logger"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache interface {
	CheckRateLimit(ctx context.Context, ip string) (bool, error)
	Ping(ctx context.Context) error
	Close() error
}

type redisCache struct {
	cache *redis.Client
}

func NewRedisCache(config config.RedisConfig) RedisCache {
	return &redisCache{
		cache: redis.NewClient(&redis.Options{
			Addr: config.RedisInternalURL,
		}),
	}
}

func (r *redisCache) CheckRateLimit(ctx context.Context, ip string) (bool, error) {
	key := "rate_limit:" + ip

	limit := int64(10)
	window := time.Duration(60) * time.Second

	count, err := r.cache.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		if err := r.cache.Expire(ctx, key, window).Err(); err != nil {
			return false, err
		}
	}

	return count <= limit, nil
}

func (r *redisCache) Ping(ctx context.Context) error {
	return r.cache.Ping(ctx).Err()
}

func (r *redisCache) Close() error {
	logger.Log.Info().Msg("closing redis cache connection")
	return r.cache.Close()
}
