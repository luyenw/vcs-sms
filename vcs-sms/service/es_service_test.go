package service

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
	"vcs-sms/model/dto"
	"vcs-sms/model/mock_entity"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCalculateUptime(t *testing.T) {
	esMock := new(mock_entity.ESMock)

	type input struct {
		startMils int64
		endMils   int64
	}
	type output struct {
		uptime []dto.ServerUptime
	}
	tests := map[string]struct {
		input
		output
		mockFunc func()
	}{
		"Test_1": {
			input: input{
				startMils: -1,
				endMils:   0,
			},
			output: output{
				uptime: []dto.ServerUptime{},
			},
			mockFunc: func() {
			},
		},
		"Test_2": {
			input: input{
				startMils: 0,
				endMils:   1715779518318,
			},
			output: output{
				uptime: []dto.ServerUptime{{}, {}},
			},
			mockFunc: func() {
				response := &dto.Response{}
				response.Aggregtions.Server.Buckets = []dto.ServerUptime{{}, {}}
				body, _ := json.Marshal(response)
				queryResponse := esapi.Response{
					Body: io.NopCloser(bytes.NewBufferString(string(body))),
				}
				esMock.On("Query", mock.Anything).Return(&queryResponse, nil)
			},
		},
	}
	service := NewESService(esMock)
	t.Run("Test_1", func(t *testing.T) {
		tests["Test_1"].mockFunc()
		uptimeInfo := service.CalculateUptime(tests["Test_1"].startMils, tests["Test_1"].endMils)
		assert.Equal(t, len(uptimeInfo), 0)
	})
	t.Run("Test_2", func(t *testing.T) {
		tests["Test_2"].mockFunc()
		uptimeInfo := service.CalculateUptime(tests["Test_2"].startMils, tests["Test_2"].endMils)
		assert.Greater(t, len(uptimeInfo), 0)
	})
}
