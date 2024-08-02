package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"vcs-sms/config"
	"vcs-sms/model/dto"
	"vcs-sms/model/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}
type IUserService interface {
	CreateNewUser(username string, password string, role entity.Role) error
	FindByUsername(username string) entity.User
	FindUserByID(id int) *entity.User
	FindRoleByName(name string) (*entity.Role, error)
	UpdateUserRole(user *entity.User, role *entity.Role) error
}

func (m *MockUserService) CreateNewUser(username string, password string, role entity.Role) error {
	args := m.Called(username, password, role)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}
func (m *MockUserService) FindByUsername(username string) entity.User {
	args := m.Called(username)
	return args.Get(0).(entity.User)
}
func (m *MockUserService) FindUserByID(id int) *entity.User {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*entity.User)
}
func (m *MockUserService) FindRoleByName(name string) (*entity.Role, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}
func (m *MockUserService) UpdateUserRole(user *entity.User, role *entity.Role) error {
	args := m.Called(user, role)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockUserService)
		controller := NewUserController(mockService)
		router := config.GetTestGin()
		router.POST("/createUser", controller.CreateUser)
		createUserRequest := dto.CreateUserRequest{
			Username: "testuser",
			Password: "password",
			Role:     entity.Role{ID: 1},
		}
		jsonValue, _ := json.Marshal(createUserRequest)
		req, _ := http.NewRequest("POST", "/createUser", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockService.On("FindByUsername", "testuser").Return(entity.User{})
		mockService.On("CreateNewUser", createUserRequest.Username, mock.Anything, createUserRequest.Role).Return(nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Username already exist", func(t *testing.T) {
		mockService := new(MockUserService)
		controller := NewUserController(mockService)
		router := config.GetTestGin()
		router.POST("/createUser", controller.CreateUser)
		createUserRequest := dto.CreateUserRequest{
			Username: "existinguser",
			Password: "password",
			Role:     entity.Role{ID: 1},
		}
		jsonValue, _ := json.Marshal(createUserRequest)
		req, _ := http.NewRequest("POST", "/createUser", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockService.On("FindByUsername", "existinguser").Return(entity.User{Username: "existinguser"})
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		mockService := new(MockUserService)
		controller := NewUserController(mockService)
		router := config.GetTestGin()
		router.POST("/createUser", controller.CreateUser)
		invalidJSON := `{"Username": "testuser", "Password": "password", "Scopes": ["read", "write"`
		req, _ := http.NewRequest("POST", "/createUser", bytes.NewBuffer([]byte(invalidJSON)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to hash password", func(t *testing.T) {
		mockService := new(MockUserService)
		controller := NewUserController(mockService)
		router := config.GetTestGin()
		router.POST("/createUser", controller.CreateUser)
		createUserRequest := dto.CreateUserRequest{
			Username: "testuser",
			Password: "passwordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpassword",
			Role:     entity.Role{ID: 1},
		}
		jsonValue, _ := json.Marshal(createUserRequest)
		req, _ := http.NewRequest("POST", "/createUser", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockService.On("FindByUsername", "testuser").Return(entity.User{})

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, `{"error":"bcrypt: password length exceeds 72 bytes"}`, w.Body.String())
	})
	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockUserService)
		controller := NewUserController(mockService)
		router := config.GetTestGin()
		router.POST("/createUser", controller.CreateUser)
		createUserRequest := dto.CreateUserRequest{
			Username: "testuser",
			Password: "password",
			Role:     entity.Role{ID: 1},
		}
		jsonValue, _ := json.Marshal(createUserRequest)
		req, _ := http.NewRequest("POST", "/createUser", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockService.On("FindByUsername", "testuser").Return(entity.User{})
		mockService.On("CreateNewUser", createUserRequest.Username, mock.Anything, createUserRequest.Role).Return(errors.New("service error"))

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// func TestUpdateUserScope(t *testing.T) {

// 	t.Run("Success", func(t *testing.T) {
// 		mockService := new(MockUserService)
// 		controller := NewUserController(mockService)
// 		router := config.GetTestGin()
// 		router.PUT("/updateUserScope/:id", controller.UpdateUserScope)
// 		updateRequest := dto.UpdatePermissionRequest{
// 			Scopes: []dto.ScopeDTO{},
// 		}
// 		jsonValue, _ := json.Marshal(updateRequest)
// 		req, _ := http.NewRequest("PUT", "/updateUserScope/1", bytes.NewBuffer(jsonValue))
// 		req.Header.Set("Content-Type", "application/json")
// 		w := httptest.NewRecorder()

// 		mockService.On("FindUserByID", 1).Return(&entity.User{ID: 1})
// 		mockService.On("UpdateUserScope", mock.Anything, []dto.ScopeDTO{}).Return(nil)

// 		router.ServeHTTP(w, req)
// 		assert.Equal(t, http.StatusOK, w.Code)
// 	})

// 	t.Run("User not found", func(t *testing.T) {
// 		mockService := new(MockUserService)
// 		controller := NewUserController(mockService)
// 		router := config.GetTestGin()
// 		router.PUT("/updateUserScope/:id", controller.UpdateUserScope)
// 		req, _ := http.NewRequest("PUT", "/updateUserScope/1", nil)
// 		req.Header.Set("Content-Type", "application/json")
// 		w := httptest.NewRecorder()
// 		mockService.On("FindUserByID", 1).Return(nil)
// 		router.ServeHTTP(w, req)
// 		assert.Equal(t, http.StatusBadRequest, w.Code)
// 	})

// 	t.Run("Invalid JSON", func(t *testing.T) {
// 		mockService := new(MockUserService)
// 		controller := NewUserController(mockService)
// 		router := config.GetTestGin()
// 		router.PUT("/updateUserScope/:id", controller.UpdateUserScope)
// 		invalidJSON := `{"Scopes": ["read", "write"`
// 		req, _ := http.NewRequest("PUT", "/updateUserScope/1", bytes.NewBuffer([]byte(invalidJSON)))
// 		req.Header.Set("Content-Type", "application/json")
// 		w := httptest.NewRecorder()

// 		mockService.On("FindUserByID", 1).Return(&entity.User{ID: 1})
// 		router.ServeHTTP(w, req)
// 		assert.Equal(t, http.StatusBadRequest, w.Code)
// 	})

// 	t.Run("Invalid ID", func(t *testing.T) {
// 		mockService := new(MockUserService)
// 		controller := NewUserController(mockService)
// 		router := config.GetTestGin()
// 		router.PUT("/updateUserScope/:id", controller.UpdateUserScope)
// 		req, _ := http.NewRequest("PUT", "/updateUserScope/invalid_id", nil)
// 		req.Header.Set("Content-Type", "application/json")
// 		w := httptest.NewRecorder()
// 		router.ServeHTTP(w, req)
// 		assert.Equal(t, http.StatusBadRequest, w.Code)
// 	})

// 	t.Run("Service Error", func(t *testing.T) {
// 		mockService := new(MockUserService)
// 		controller := NewUserController(mockService)
// 		router := config.GetTestGin()
// 		router.PUT("/updateUserScope/:id", controller.UpdateUserScope)
// 		updateRequest := dto.UpdatePermissionRequest{
// 			Scopes: []dto.ScopeDTO{},
// 		}
// 		jsonValue, _ := json.Marshal(updateRequest)
// 		req, _ := http.NewRequest("PUT", "/updateUserScope/1", bytes.NewBuffer(jsonValue))
// 		req.Header.Set("Content-Type", "application/json")
// 		w := httptest.NewRecorder()

// 		mockService.On("FindUserByID", 1).Return(&entity.User{ID: 1})
// 		mockService.On("UpdateUserScope", &entity.User{ID: 1}, []dto.ScopeDTO{}).Return(errors.New("service error"))

// 		router.ServeHTTP(w, req)
// 		assert.Equal(t, http.StatusInternalServerError, w.Code)
// 	})
// }
