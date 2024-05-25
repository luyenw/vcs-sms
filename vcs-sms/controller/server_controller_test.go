package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"vcs-sms/config"
	"vcs-sms/model/dto"
	"vcs-sms/model/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockServerService struct {
	mock.Mock
}
type MockCacheService struct {
	mock.Mock
}
type MockXLSXService struct {
	mock.Mock
}

func (m *MockServerService) CreateServer(server *entity.Server) error {
	args := m.Called(server)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}
func (m *MockServerService) DeleteServerById(id int) error {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}
func (m *MockServerService) FindServerById(id int) *entity.Server {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*entity.Server)
}
func (m *MockServerService) GetAllServers() []entity.Server {
	args := m.Called()
	return args.Get(0).([]entity.Server)
}
func (m *MockServerService) GetServer(queryParam *dto.QueryParam) []entity.Server {
	args := m.Called(queryParam)
	return args.Get(0).([]entity.Server)
}
func (m *MockServerService) UpdateServer(server *entity.Server) error {
	args := m.Called(server)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

func (m *MockCacheService) Get(key string) (string, error) {
	args := m.Called()
	if args.Get(1) == nil {
		return "", nil
	}
	return args.String(0), args.Error(1)
}
func (m *MockCacheService) Set(key string, value interface{}) error {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}
func (m *MockXLSXService) ExportXLSX(servers []entity.Server) (string, error) {
	args := m.Called(servers)
	if args.Get(1) == nil {
		return "", nil
	}
	return args.String(0), args.Error(1)
}
func (m *MockXLSXService) ImportXLSX(filePath string) ([][]string, error) {
	args := m.Called(filePath)
	if args.Get(1) == nil {
		return args.Get(0).([][]string), nil
	}
	return args.Get(0).([][]string), args.Error(1)
}
func TestGetServer(t *testing.T) {
	t.Run("Failed to bind query param", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.GET("/getServer", serverController.GetServer)

		req, _ := http.NewRequest("GET", "/getServer?page=abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Success", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.GET("/getServer", serverController.GetServer)
		req, _ := http.NewRequest("GET", "/getServer?page=1", nil)

		mockServerService.On("GetServer", &dto.QueryParam{Page: 1}).Return([]entity.Server{})
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
func TestUpdateServer(t *testing.T) {
	t.Run("Failed to parsing id", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.PUT("/updateServer/:id", serverController.UpdateServer)
		req, _ := http.NewRequest("PUT", "/updateServer/abc", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Failed to find server by id", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.PUT("/updateServer/:id", serverController.UpdateServer)
		req, _ := http.NewRequest("PUT", "/updateServer/1", nil)

		mockServerService.On("FindServerById", 1).Return(nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("Failed to bind input server", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.PUT("/updateServer/:id", serverController.UpdateServer)
		jsonValue := `{"name": "server1", "ipv4": "`
		req, _ := http.NewRequest("PUT", "/updateServer/1", bytes.NewBuffer([]byte(jsonValue)))

		mockServerService.On("FindServerById", 1).Return(&entity.Server{})
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Failed to update server", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.PUT("/updateServer/:id", serverController.UpdateServer)
		inputRequest := dto.InputServer{
			Name: "server1",
		}
		jsonValue, _ := json.Marshal(inputRequest)
		req, _ := http.NewRequest("PUT", "/updateServer/1", bytes.NewBuffer([]byte(jsonValue)))
		mockServerService.On("FindServerById", 1).Return(&entity.Server{})
		mockServerService.On("UpdateServer", mock.Anything).Return(errors.New("error"))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("Failed on set cache", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.PUT("/updateServer/:id", serverController.UpdateServer)
		inputRequest := dto.InputServer{
			Name: "server1",
		}
		jsonValue, _ := json.Marshal(inputRequest)
		req, _ := http.NewRequest("PUT", "/updateServer/1", bytes.NewBuffer([]byte(jsonValue)))
		mockServerService.On("FindServerById", 1).Return(&entity.Server{})
		mockServerService.On("UpdateServer", mock.Anything).Return(nil)
		mockCacheService.On("Set", mock.Anything, mock.Anything).Return(errors.New("error"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("Success", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.PUT("/updateServer/:id", serverController.UpdateServer)
		inputRequest := dto.InputServer{
			Name: "server1",
		}
		jsonValue, _ := json.Marshal(inputRequest)
		req, _ := http.NewRequest("PUT", "/updateServer/1", bytes.NewBuffer([]byte(jsonValue)))
		mockServerService.On("FindServerById", 1).Return(&entity.Server{})
		mockServerService.On("UpdateServer", mock.Anything).Return(nil)
		mockCacheService.On("Set", mock.Anything, mock.Anything).Return(nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestDeleteServer(t *testing.T) {
	t.Run("Failed to parsing id", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.DELETE("/deleteServer/:id", serverController.DeleteServer)
		req, _ := http.NewRequest("DELETE", "/deleteServer/abc", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Service error", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.DELETE("/deleteServer/:id", serverController.DeleteServer)
		req, _ := http.NewRequest("DELETE", "/deleteServer/1", nil)

		mockServerService.On("DeleteServerById", 1).Return(errors.New("error"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("Failed to set cache", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.DELETE("/deleteServer/:id", serverController.DeleteServer)
		req, _ := http.NewRequest("DELETE", "/deleteServer/1", nil)

		mockServerService.On("DeleteServerById", 1).Return(nil)
		mockCacheService.On("Set", mock.Anything, mock.Anything).Return(errors.New("error"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("Success", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.DELETE("/deleteServer/:id", serverController.DeleteServer)
		req, _ := http.NewRequest("DELETE", "/deleteServer/1", nil)

		mockServerService.On("DeleteServerById", 1).Return(nil)
		mockCacheService.On("Set", mock.Anything, mock.Anything).Return(nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
func TestCreateServer(t *testing.T) {
	t.Run("Failed to bind json", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.POST("/createServer", serverController.CreateServer)
		jsonValue := `{"name": "server1`
		req, _ := http.NewRequest("POST", "/createServer", bytes.NewBuffer([]byte(jsonValue)))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Failed on data validation", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.POST("/createServer", serverController.CreateServer)

		server := dto.InputServer{
			Name:   "server",
			IPv4:   "notipv4format",
			Status: 4,
		}
		jsonValue, _ := json.Marshal(server)
		req, _ := http.NewRequest("POST", "/createServer", bytes.NewBuffer([]byte(jsonValue)))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Service error", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.POST("/createServer", serverController.CreateServer)
		server := dto.InputServer{
			Name:   "server",
			IPv4:   "notipv4format",
			Status: 4,
		}
		jsonValue, _ := json.Marshal(server)
		req, _ := http.NewRequest("POST", "/createServer", bytes.NewBuffer([]byte(jsonValue)))

		mockServerService.On("CreateServer", mock.Anything).Return(errors.New("error"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Failed to set cache", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.POST("/createServer", serverController.CreateServer)
		server := dto.InputServer{
			Name:   "server",
			IPv4:   "0.0.0.0",
			Status: 1,
		}
		jsonValue, _ := json.Marshal(server)
		req, _ := http.NewRequest("POST", "/createServer", bytes.NewBuffer([]byte(jsonValue)))

		mockServerService.On("CreateServer", mock.Anything).Return(nil)
		mockServerService.On("GetAllServers").Return([]entity.Server{})
		mockCacheService.On("Set", mock.Anything, mock.Anything).Return(errors.New("error"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("Success", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.POST("/createServer", serverController.CreateServer)
		server := dto.InputServer{
			Name:   "server",
			IPv4:   "0.0.0.0",
			Status: 1,
		}
		jsonValue, _ := json.Marshal(server)
		req, _ := http.NewRequest("POST", "/createServer", bytes.NewBuffer([]byte(jsonValue)))

		mockServerService.On("CreateServer", mock.Anything).Return(nil)
		mockServerService.On("GetAllServers").Return([]entity.Server{})
		mockCacheService.On("Set", mock.Anything, mock.Anything).Return(nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
func TestExportServers(t *testing.T) {
	t.Run("Failed to bind query param", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.GET("/exportServer", serverController.ExportServers)
		req, _ := http.NewRequest("GET", "/exportServer?page=abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Failed to export xlsx", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.GET("/exportServer", serverController.ExportServers)
		req, _ := http.NewRequest("GET", "/exportServer?page=1", nil)
		mockServerService.On("GetServer", &dto.QueryParam{Page: 1}).Return([]entity.Server{})
		mockXLSXService.On("ExportXLSX", mock.Anything).Return("", errors.New("error"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("Success", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.GET("/exportServer", serverController.ExportServers)
		req, _ := http.NewRequest("GET", "/exportServer?page_size=10&page=1", nil)
		mockServerService.On("GetServer", &dto.QueryParam{Page: 1, PageSize: 10}).Return([]entity.Server{})
		mockXLSXService.On("ExportXLSX", mock.Anything).Return("url", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestImportServers(t *testing.T) {
	t.Run("Failed to get file", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.POST("/importServer", serverController.ImportServers)
		form := map[string]string{"name": "John"}
		ct, body, _ := createForm(form)
		req, _ := http.NewRequest("POST", "/importServer", body)
		req.Header.Set("Content-Type", ct)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"error":"Cannot get file from request"}`, w.Body.String())
	})
	t.Run("Invalid file format", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.POST("/importServer", serverController.ImportServers)
		form := map[string]string{"input": "@/mnt/c/Users/luyen/Desktop/vcs-sms/vcs-sms/go.mod"}
		ct, body, _ := createForm(form)
		req, _ := http.NewRequest("POST", "/importServer", body)
		req.Header.Set("Content-Type", ct)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"error":"Invalid file format"}`, w.Body.String())
	})
	t.Run("Failed to import file", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.POST("/importServer", serverController.ImportServers)
		form := map[string]string{"input": "@/mnt/c/Users/luyen/Desktop/vcs-sms/vcs-sms/tmp/template.xlsx"}
		ct, body, _ := createForm(form)
		req, _ := http.NewRequest("POST", "/importServer", body)
		req.Header.Set("Content-Type", ct)

		mockXLSXService.On("ImportXLSX", mock.Anything).Return([][]string{}, errors.New("error"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, `{"error":"Failed to import file"}`, w.Body.String())
	})
	t.Run("Failed to set cache", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.POST("/importServer", serverController.ImportServers)
		form := map[string]string{"input": "@/mnt/c/Users/luyen/Desktop/vcs-sms/vcs-sms/tmp/template.xlsx"}
		ct, body, _ := createForm(form)
		req, _ := http.NewRequest("POST", "/importServer", body)
		req.Header.Set("Content-Type", ct)

		mockXLSXService.On("ImportXLSX", mock.Anything).Return([][]string{
			{"name", "ipv4"},
			{"server1", "0.0.0.0"},
		}, nil)
		mockServerService.On("CreateServer", mock.Anything).Return(nil)
		mockServerService.On("GetAllServers").Return([]entity.Server{})
		mockCacheService.On("Set", mock.Anything, mock.Anything).Return(errors.New("error"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("Failed to parse server", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.POST("/importServer", serverController.ImportServers)
		form := map[string]string{"input": "@/mnt/c/Users/luyen/Desktop/vcs-sms/vcs-sms/tmp/template.xlsx"}
		ct, body, _ := createForm(form)
		req, _ := http.NewRequest("POST", "/importServer", body)
		req.Header.Set("Content-Type", ct)

		mockXLSXService.On("ImportXLSX", mock.Anything).Return([][]string{
			{"name", "ipv4"},
			{"server1", "0.0.0abc"},
		}, nil)
		mockServerService.On("CreateServer", mock.Anything).Return(nil)
		mockServerService.On("GetAllServers").Return([]entity.Server{})
		mockCacheService.On("Set", mock.Anything, mock.Anything).Return(errors.New("error"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "{\"data\":{\"failure_count\":1,\"success_count\":0,\"success_names\":[]},\"message\":\"Imported successfully\"}", w.Body.String())
	})
	t.Run("Service error", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.POST("/importServer", serverController.ImportServers)
		form := map[string]string{"input": "@/mnt/c/Users/luyen/Desktop/vcs-sms/vcs-sms/tmp/template.xlsx"}
		ct, body, _ := createForm(form)
		req, _ := http.NewRequest("POST", "/importServer", body)
		req.Header.Set("Content-Type", ct)

		mockXLSXService.On("ImportXLSX", mock.Anything).Return([][]string{
			{"name", "ipv4"},
			{"server1", "0.0.0.0"},
		}, nil)
		mockServerService.On("CreateServer", mock.Anything).Return(errors.New("error"))
		mockServerService.On("GetAllServers").Return([]entity.Server{})
		mockCacheService.On("Set", mock.Anything, mock.Anything).Return(errors.New("error"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "{\"data\":{\"failure_count\":1,\"success_count\":0,\"success_names\":[]},\"message\":\"Imported successfully\"}", w.Body.String())
	})
	t.Run("Success", func(t *testing.T) {
		mockServerService := new(MockServerService)
		mockCacheService := new(MockCacheService)
		mockXLSXService := new(MockXLSXService)
		serverController := NewServerController(mockServerService, mockCacheService, mockXLSXService)
		router := config.GetTestGin()
		router.POST("/importServer", serverController.ImportServers)
		form := map[string]string{"input": "@/mnt/c/Users/luyen/Desktop/vcs-sms/vcs-sms/tmp/template.xlsx"}
		ct, body, _ := createForm(form)
		req, _ := http.NewRequest("POST", "/importServer", body)
		req.Header.Set("Content-Type", ct)

		mockXLSXService.On("ImportXLSX", mock.Anything).Return([][]string{
			{"name", "ipv4"},
			{"server1", "0.0.0.0"},
		}, nil)
		mockServerService.On("CreateServer", mock.Anything).Return(nil)
		mockServerService.On("GetAllServers").Return([]entity.Server{})
		mockCacheService.On("Set", mock.Anything, mock.Anything).Return(nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
func createForm(form map[string]string) (string, io.Reader, error) {
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	defer mp.Close()
	for key, val := range form {
		if strings.HasPrefix(val, "@") {
			val = val[1:]
			file, err := os.Open(val)
			if err != nil {
				return "", nil, err
			}
			defer file.Close()
			part, err := mp.CreateFormFile(key, val)
			if err != nil {
				return "", nil, err
			}
			io.Copy(part, file)
		} else {
			mp.WriteField(key, val)
		}
	}
	return mp.FormDataContentType(), body, nil
}
