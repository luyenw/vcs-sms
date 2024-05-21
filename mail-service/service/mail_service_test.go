package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMail(t *testing.T) {
	type input struct {
		to []string
	}
	type output struct {
		err error
	}
	tests := map[string]struct {
		input
		output
	}{
		"Test_1": {
			input: input{
				to: []string{
					"test1@gmaill.coml",
				},
			},
			output: output{
				err: nil,
			},
		},
		"Test_2": {
			input: input{
				to: []string{
					"test1@gmail.com\n",
				},
			},
			output: output{
				err: errors.New("smtp: A line must not contain CR or LF"),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			service := NewMailService()
			err := service.SendEmail(tc.input.to, "test")
			assert.Equal(t, tc.output.err, err)
		})
	}
}
