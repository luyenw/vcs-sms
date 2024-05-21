package service

// import (
// 	"testing"
// 	"vcs-sms/model/entity"
// 	"vcs-sms/model/mock_entity"
// )

// func TestFindByUsername(t *testing.T) {
// 	mockDB := mock_entity.NewMockDatabase()
// 	// service := NewUserService(mockDB)
// 	type input struct {
// 		username string
// 	}
// 	type expectation struct {
// 		user entity.User
// 	}
// 	tests := map[string]struct {
// 		input
// 		expectation
// 		mockFunc func()
// 	}{
// 		"Test_1": {
// 			input: input{
// 				username: "username",
// 			},
// 			expectation: expectation{
// 				user: entity.User{
// 					Username: "username",
// 				},
// 			},
// 			mockFunc: func() {
// 				mockDB.On("Where", "username = ?", []interface{}{"username"}).Return(nil)
// 				mockDB.On("Preload", "Scopes").Return(nil)
// 				mockDB.On("First", &entity.User{}).Return(&entity.User{Username: "username"}, nil)
// 			},
// 		},
// 	}
// 	for name, test := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			test.mockFunc()
// 			// user := service.FindByUsername(test.input.username)
// 			// if user.Username != test.expectation.user.Username {
// 			// 	t.Errorf("Expected: %v but got: %v", test.expectation.user.Username, user.Username)
// 			// }
// 		})
// 	}
// }
