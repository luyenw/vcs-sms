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
	"golang.org/x/crypto/bcrypt"
)

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(user *entity.User) (string, error) {
	args := m.Called(user)
	if args.Get(1) == nil {
		return args.String(0), nil
	}
	return args.String(0), args.Error(1)
}

func TestLogin(t *testing.T) {
	t.Run("Failed to bind JSON", func(t *testing.T) {
		jwtService := new(MockJWTService)
		userService := new(MockUserService)
		controller := NewAuthController(userService, jwtService)
		router := config.GetTestGin()
		router.POST("/login", controller.Login)

		loginRequest := `{"username": "test", "password": "test"`
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(loginRequest)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to find user", func(t *testing.T) {
		jwtService := new(MockJWTService)
		userService := new(MockUserService)
		controller := NewAuthController(userService, jwtService)
		router := config.GetTestGin()
		router.POST("/login", controller.Login)

		loginRequest := dto.AuthRequest{
			Username: "test",
			Password: "test",
		}
		jsonValue, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(jsonValue)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		userService.On("FindByUsername", "test").Return(entity.User{})
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Failed to compare password", func(t *testing.T) {
		jwtService := new(MockJWTService)
		userService := new(MockUserService)
		controller := NewAuthController(userService, jwtService)
		router := config.GetTestGin()
		router.POST("/login", controller.Login)

		loginRequest := dto.AuthRequest{
			Username: "test",
			Password: "test",
		}
		jsonValue, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(jsonValue)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		userService.On("FindByUsername", "test").Return(entity.User{
			Username: "test",
			Password: "different_hash_password",
		})
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Failed on generate token", func(t *testing.T) {
		jwtService := new(MockJWTService)
		userService := new(MockUserService)
		controller := NewAuthController(userService, jwtService)
		router := config.GetTestGin()
		router.POST("/login", controller.Login)

		loginRequest := dto.AuthRequest{
			Username: "test",
			Password: "test",
		}
		jsonValue, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(jsonValue)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		hashPwd, _ := bcrypt.GenerateFromPassword([]byte(loginRequest.Password), bcrypt.DefaultCost)
		userService.On("FindByUsername", "test").Return(entity.User{
			Username: "test",
			Password: string(hashPwd),
		})
		jwtService.On("GenerateToken", mock.Anything).Return("", errors.New("error"))
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Success", func(t *testing.T) {
		jwtService := new(MockJWTService)
		userService := new(MockUserService)
		controller := NewAuthController(userService, jwtService)
		router := config.GetTestGin()
		router.POST("/login", controller.Login)

		loginRequest := dto.AuthRequest{
			Username: "test",
			Password: "test",
		}
		jsonValue, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(jsonValue)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		hashPwd, _ := bcrypt.GenerateFromPassword([]byte(loginRequest.Password), bcrypt.DefaultCost)
		userService.On("FindByUsername", "test").Return(entity.User{
			Username: "test",
			Password: string(hashPwd),
		})
		jwtService.On("GenerateToken", mock.Anything).Return("token", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRegistere(t *testing.T) {
	t.Run("Failed to bind JSON", func(t *testing.T) {
		jwtService := new(MockJWTService)
		userService := new(MockUserService)
		controller := NewAuthController(userService, jwtService)
		router := config.GetTestGin()
		router.POST("/register", controller.Register)
		jsonValue := `{"username": "test", "password": "test"`
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(jsonValue)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Username already exist", func(t *testing.T) {
		jwtService := new(MockJWTService)
		userService := new(MockUserService)
		controller := NewAuthController(userService, jwtService)
		router := config.GetTestGin()
		router.POST("/register", controller.Register)
		authRequest := dto.AuthRequest{
			Username: "test",
			Password: "test",
		}
		jsonValue, _ := json.Marshal(authRequest)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(jsonValue)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		userService.On("FindByUsername", "test").Return(entity.User{Username: "test"})
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to hash password", func(t *testing.T) {
		jwtService := new(MockJWTService)
		userService := new(MockUserService)
		controller := NewAuthController(userService, jwtService)
		router := config.GetTestGin()
		router.POST("/register", controller.Register)
		authRequest := dto.AuthRequest{
			Username: "test",
			Password: "longerthan72byteslongerthan72byteslongerthan72byteslongerthan72byteslongerthan72bytes",
		}
		jsonValue, _ := json.Marshal(authRequest)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(jsonValue)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		userService.On("FindByUsername", "test").Return(entity.User{})
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Service error", func(t *testing.T) {
		jwtService := new(MockJWTService)
		userService := new(MockUserService)
		controller := NewAuthController(userService, jwtService)
		router := config.GetTestGin()
		router.POST("/register", controller.Register)
		authRequest := dto.AuthRequest{
			Username: "test",
			Password: "test",
		}
		jsonValue, _ := json.Marshal(authRequest)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(jsonValue)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		userService.On("FindByUsername", "test").Return(entity.User{})
		userService.On("CreateNewUser", authRequest.Username, mock.Anything, []entity.Scope{}).Return(nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Success", func(t *testing.T) {
		jwtService := new(MockJWTService)
		userService := new(MockUserService)
		controller := NewAuthController(userService, jwtService)
		router := config.GetTestGin()
		router.POST("/register", controller.Register)
		authRequest := dto.AuthRequest{
			Username: "test",
			Password: "test",
		}
		jsonValue, _ := json.Marshal(authRequest)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(jsonValue)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		userService.On("FindByUsername", "test").Return(entity.User{})
		userService.On("CreateNewUser", authRequest.Username, mock.Anything, []entity.Scope{}).Return(nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
