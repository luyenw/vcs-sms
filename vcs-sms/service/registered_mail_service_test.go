package service

import (
	"testing"
	"vcs-sms/model/entity"
	"vcs-sms/model/mock_entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllRegisteredMails(t *testing.T) {
	t.Run("Get All", func(t *testing.T) {
		mockDB := mock_entity.NewMockDatabase()
		mockDB.On("Find", &[]entity.RegisteredEmail{}, mock.Anything).Return(nil, []entity.RegisteredEmail{})
		service := NewRegisteredMailService(mockDB)
		mails := service.GetAllRegisteredMails()
		assert.NotNil(t, mails)
	})
}
