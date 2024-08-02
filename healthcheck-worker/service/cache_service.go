package service

import (
	"context"
	"fmt"
	"healthcheck-worker/config/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
}
type CacheService struct {
	cache CacheRepository
}

func NewCacheService(cache CacheRepository) *CacheService {
	return &CacheService{
		cache: cache,
	}
}

func (cs *CacheService) Set(key string, value interface{}) error {
	log := logger.NewLogger()
	log.Info(fmt.Sprintf("Redis Set key: %s - value: %s", key, value))
	err := cs.cache.Set(context.Background(), key, value, 0).Err()
	return err
}

func (cs *CacheService) Get(key string) (string, error) {
	val, err := cs.cache.Get(context.Background(), key).Result()
	return val, err
}
