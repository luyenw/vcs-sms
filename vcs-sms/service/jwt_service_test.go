package service

import (
	"testing"
	"vcs-sms/model/entity"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	type input struct {
		user *entity.User
	}
	type expectation struct {
		token string
		err   error
	}
	tests := map[string]struct {
		input
		expectation
	}{
		"Test_1": {
			input: input{
				user: &entity.User{
					Username: "luyendd",
					Scopes: []entity.Scope{
						{Name: "all"},
					},
				},
			},
			expectation: expectation{
				token: "aasASa",
				err:   nil,
			},
		},
	}
	t.Run("Test_1", func(t *testing.T) {
		service := NewJWTService()
		token, err := service.GenerateToken(tests["Test_1"].user)
		assert.Greater(t, len(token), 0)
		assert.Equal(t, tests["Test_1"].expectation.err, err)
	})

}
