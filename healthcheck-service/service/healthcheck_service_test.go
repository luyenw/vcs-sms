package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	type input struct {
		addr string
	}
	type expectation struct {
		result int
	}
	tests := map[string]struct {
		input
		expectation
	}{
		"Test_1": {
			input: input{
				addr: "1.1.1.1",
			},
			expectation: expectation{
				result: 1,
			},
		},
		"Test_2": {
			input: input{
				addr: "0.0.0.0",
			},
			expectation: expectation{
				result: 0,
			},
		},
		"Test_3": {
			input: input{
				addr: "example",
			},
			expectation: expectation{
				result: 0,
			},
		},
		"Test_4": {
			input: input{
				addr: "google.com",
			},
			expectation: expectation{
				result: 1,
			},
		},
	}
	t.Run("Test_1", func(t *testing.T) {
		test := tests["Test_1"]
		service := NewHealthCheckService()
		result := service.HealthCheck(test.addr)
		assert.Equal(t, test.expectation.result, result)
	})
}
