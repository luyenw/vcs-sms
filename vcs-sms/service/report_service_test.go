package service

// import (
// 	"errors"
// 	"testing"
// 	"vcs-sms/model/mock_entity"
// 	"vcs-sms/repo"

// 	"github.com/stretchr/testify/assert"
// )

// func TestSendReport(t *testing.T) {
// 	type input struct {
// 		startMilis int64
// 		endMilis   int64
// 		to         []string
// 	}
// 	type output struct {
// 		err error
// 	}
// 	tests := map[string]struct {
// 		input
// 		output
// 	}{
// 		"Test_1": {
// 			input: input{
// 				startMilis: 0,
// 				endMilis:   -1,
// 				to:         []string{},
// 			},
// 			output: output{
// 				err: errors.New("Invalid time range"),
// 			},
// 		},
// 		"Test_2": {
// 			input: input{
// 				startMilis: 0,
// 				endMilis:   1715779518318,
// 				to:         []string{"tudu8603@gail.com"},
// 			},
// 			output: output{
// 				err: nil,
// 			},
// 		},
// 		"Test_3": {
// 			input: input{
// 				startMilis: 0,
// 				endMilis:   1715779518318,
// 				to:         []string{"tudu86032gail.com\n"},
// 			},
// 			output: output{
// 				err: errors.New("smtp: A line must not contain CR or LF"),
// 			},
// 		},
// 	}

// 	service := NewReportService(NewESService(&repo.ESClient{Client: mock_entity.NewESMock()}),
// 		NewRegisteredMailService(mock_entity.NewMockDatabase()),
// 		NewServerService(mock_entity.NewMockDatabase()),
// 		NewCacheService(mock_entity.NewMockRedis()),
// 	)

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			err := service.SendReport(tc.startMilis, tc.endMilis, tc.to)
// 			assert.Equal(t, tc.err, err)
// 		})
// 	}

// }
