package service

import (
	"testing"
	"vcs-sms/model/entity"

	"github.com/stretchr/testify/assert"
)

func TestGetAllRegisteredMails(t *testing.T) {
	type input struct {
	}
	type output struct {
		mails []entity.RegisteredEmail
	}
	tests := map[string]struct {
		input
		output
	}{
		"Test_1": {
			input: input{},
			output: output{
				mails: []entity.RegisteredEmail{},
			},
		},
	}
	for name := range tests {
		t.Run(name, func(t *testing.T) {
			service := NewRegisteredMailService()
			mails := service.GetAllRegisteredMails()
			assert.NotNil(t, mails)
		})
	}
}
