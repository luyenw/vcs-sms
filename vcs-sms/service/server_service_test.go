package service

import (
	"errors"
	"testing"
	"vcs-sms/model/entity"
	"vcs-sms/model/mock_entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestGetAllServers(t *testing.T) {
	mockDB := mock_entity.NewMockDatabase()
	mockDB_2 := mock_entity.NewMockDatabase()

	var service *ServerService

	type expectation struct {
		servers []entity.Server
	}

	tests := map[string]struct {
		expectation
		mockFunc func()
	}{
		"Test_1": {
			expectation: expectation{
				servers: []entity.Server{
					{ID: 1},
				},
			},
			mockFunc: func() {
				mockDB.On("Find", &[]entity.Server{}, mock.Anything).Return(nil, &[]entity.Server{{ID: 1}})
				service = NewServerService(mockDB)
			},
		},
		"Test_2": {
			expectation: expectation{
				servers: []entity.Server{},
			},
			mockFunc: func() {
				mockDB_2.On("Find", &[]entity.Server{}, mock.Anything).Return(errors.New("error"), nil)
				service = NewServerService(mockDB_2)
			},
		},
	}

	t.Run("Test_1", func(t *testing.T) {
		tests["Test_1"].mockFunc()
		servers := service.GetAllServers()
		assert.Equal(t, len(tests["Test_1"].expectation.servers), len(servers))
	})

	t.Run("Test_2", func(t *testing.T) {
		tests["Test_2"].mockFunc()
		servers := service.GetAllServers()
		assert.Equal(t, len(tests["Test_2"].expectation.servers), len(servers))
	})
}

func TestDeleteServer(t *testing.T) {
	mockDB := mock_entity.NewMockDatabase()
	service := NewServerService(mockDB)
	type input struct {
		id int
	}
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		input
		expectation
		mockFunc func()
	}{
		"Pass": {
			input: input{
				id: 1,
			},
			expectation: expectation{
				err: nil,
			},
			mockFunc: func() {
				mockDB.On("Delete").Return(nil)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mockFunc()
			err := service.DeleteServerById(test.input.id)
			assert.Equal(t, test.expectation.err, err)
		})
	}
}

func TestUpdateServer(t *testing.T) {
	mockDB := mock_entity.NewMockDatabase()
	service := NewServerService(mockDB)
	type input struct {
		server *entity.Server
	}
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		input
		expectation
		mockFunc func()
	}{
		"Pass": {
			input: input{
				server: &entity.Server{},
			},
			expectation: expectation{
				err: nil,
			},
			mockFunc: func() {
				mockDB.On("Save").Return(nil)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mockFunc()
			err := service.UpdateServer(test.input.server)
			assert.Equal(t, test.expectation.err, err)
		})
	}
}

func TestFindServerById(t *testing.T) {
	mockDB := mock_entity.NewMockDatabase()
	service := NewServerService(mockDB)
	type input struct {
		id int
	}
	type expectation struct {
		server *entity.Server
	}
	tests := map[string]struct {
		input
		expectation
		mockFunc func()
	}{
		"Pass": {
			input: input{
				id: 1,
			},
			expectation: expectation{
				server: &entity.Server{
					ID: 1,
				},
			},
			mockFunc: func() {
				mockDB.On("First", &entity.Server{}, []interface{}{1}).Return(&entity.Server{ID: 1}, nil)
			},
		},
		"Fail": {
			input: input{
				id: 2,
			},
			expectation: expectation{
				server: nil,
			},
			mockFunc: func() {
				mockDB.On("First", &entity.Server{}, []interface{}{2}).Return(nil, gorm.ErrRecordNotFound)
			},
		},
	}
	t.Run("Pass", func(t *testing.T) {
		tests["Pass"].mockFunc()
		server := service.FindServerById(tests["Pass"].input.id)
		assert.Equal(t, tests["Pass"].expectation.server.ID, server.ID)
	})
	t.Run("Fail", func(t *testing.T) {
		tests["Fail"].mockFunc()
		server := service.FindServerById(tests["Fail"].input.id)
		assert.Nil(t, server)
	})
}

func TestCreateServer(t *testing.T) {
	mock := mock_entity.NewMockDatabase()
	mock_2 := mock_entity.NewMockDatabase()
	var service *ServerService

	type input struct {
		server *entity.Server
	}
	type expectation struct {
		err error
	}
	tests := map[string]struct {
		input
		expectation
		mockFunc func()
	}{
		"Test 1": {
			input: input{
				server: &entity.Server{},
			},
			expectation: expectation{
				err: nil,
			},
			mockFunc: func() {
				mock.On("Create").Return(nil)
				service = NewServerService(mock)
			},
		},
		"Test 2": {
			input: input{
				server: &entity.Server{},
			},
			expectation: expectation{
				err: errors.New("error"),
			},
			mockFunc: func() {
				mock_2.On("Create").Return(errors.New("error"))
				service = NewServerService(mock_2)
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mockFunc()
			err := service.CreateServer(test.input.server)
			if err != nil {
				if err.Error() != test.expectation.err.Error() {
					t.Errorf("Expected: %v but got: %v", test.expectation.err, err)
				}
			} else {
				if test.expectation.err != nil {
					t.Errorf("Expected: %v but got: %v", test.expectation.err, err)
				}
			}
		})
	}
}

// func TestGetServer(t *testing.T) {
// 	mockDB := mock_entity.NewMockDatabase()
// 	// mock_2 := mock_entity.NewMockDatabase()
// 	var service *ServerService

// 	type input struct {
// 		query *dto.QueryParam
// 	}
// 	type expectation struct {
// 		servers []entity.Server
// 	}
// 	tests := map[string]struct {
// 		input
// 		expectation
// 		mockFunc func()
// 	}{
// 		"Test_1": {
// 			input: input{
// 				query: &dto.QueryParam{},
// 			},
// 			expectation: expectation{
// 				servers: []entity.Server{
// 					{ID: 1},
// 				},
// 			},
// 			mockFunc: func() {
// 				mockDB.On("Order", "created_time desc").Return(nil)
// 				mockDB.On("Offset", 10).Return(nil)
// 				mockDB.On("Limit", 10).Return(nil)
// 				mockDB.On("Find", &[]entity.Server{}, mock.Anything).Return(nil, &[]entity.Server{{ID: 1}})
// 				service = NewServerService(mockDB)
// 			},
// 		},
// 	}
// 	t.Run("Test_1", func(t *testing.T) {
// 		tests["Test_1"].mockFunc()
// 		servers := service.GetServer(tests["Test_1"].input.query)
// 		assert.Equal(t, len(tests["Test_1"].expectation.servers), len(servers))
// 	})
// }
