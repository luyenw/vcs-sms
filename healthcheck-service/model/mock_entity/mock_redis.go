package mock_entity

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

type MockRedis struct {
	mock.Mock
}

func NewMockRedis() *MockRedis {
	return &MockRedis{}
}

func (m *MockRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called()
	return args.Get(0).(*redis.StatusCmd)
}
func (m *MockRedis) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called()
	cli := redis.NewStringCmd(context.Background(), nil)
	cli.SetVal(args.String(0))
	return cli
}
