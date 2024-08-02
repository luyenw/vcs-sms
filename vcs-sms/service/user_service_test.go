package service

import (
	"testing"
	"vcs-sms/model/entity"
	"vcs-sms/model/mock_entity"

	"github.com/stretchr/testify/mock"
)

func TestFindByUsername(t *testing.T) {
	t.Run("Test1", func(t *testing.T) {
		mockDB := mock_entity.NewMockDatabase()
		service := NewUserService(mockDB)

		mockDB.On("Where", "username = ?", []interface{}{"username"}).Return(nil)
		mockDB.On("Preload", mock.Anything, []interface{}{mock.Anything}).Return(nil)
		mockDB.On("First", &entity.User{}).Return(&entity.User{Username: "username"}, nil)
		user := service.FindByUsername("username")
		if user.Username != "username" {
			t.Errorf("Expected: %v but got: %v", "username", user.Username)
		}
	})
}
