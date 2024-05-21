package service

import (
	"context"
	"testing"
	"vcs-sms/model/mock_entity"

	"github.com/redis/go-redis/v9"
)

func TestCacheSet(t *testing.T) {
	mockCache := mock_entity.NewMockRedis()
	service := NewCacheService(mockCache)
	type input struct {
		key   string
		value interface{}
	}
	type expectation struct {
		err   error
		value string
	}
	tests := map[string]struct {
		input
		expectation
		mockFuc func()
	}{
		"Pass": {
			input: input{
				key:   "key",
				value: "value",
			},
			expectation: expectation{
				err: nil,
			},
			mockFuc: func() {
				mockCache.On("Set").Return(redis.NewStatusCmd(context.Background(), nil))
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mockFuc()
			err := service.Set(test.input.key, test.input.value)
			if err != test.expectation.err {
				t.Errorf("Expected: %v but got: %v", test.expectation.err, err)
			}
		})
	}
}

func TestCacheGet(t *testing.T) {
	mockDB := mock_entity.NewMockRedis()
	service := NewCacheService(mockDB)

	type input struct {
		key string
	}
	type expectation struct {
		value string
		err   error
	}
	tests := map[string]struct {
		input
		expectation
		mockFuc func()
	}{
		"Test_1": {
			input: input{
				key: "key",
			},

			expectation: expectation{
				value: "value",
				err:   nil,
			},
			mockFuc: func() {
				mockDB.On("Get").Return("value", nil)
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mockFuc()
			val, err := service.Get(test.input.key)
			if err != test.expectation.err {
				t.Errorf("Expected: %v but got: %v", test.expectation.err, err)
			}
			if val != test.expectation.value {
				t.Errorf("Expected: %v but got: %v", test.expectation.value, val)
			}
		})
	}

}
