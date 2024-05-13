package service

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	cli *redis.Client
}

func NewCacheService(cli *redis.Client) *CacheService {
	return &CacheService{
		cli: cli,
	}
}

func (cs *CacheService) Set(key string, value interface{}) error {
	log.Info(fmt.Sprintf("Redis Set key: %s - value: %s", key, value))
	err := cs.cli.Set(context.Background(), key, value, 0).Err()
	return err
}

func (cs *CacheService) Get(key string) (string, error) {
	log.Info(fmt.Sprintf("Redis Get key: %s", key))
	val, err := cs.cli.Get(context.Background(), key).Result()
	return val, err
}
