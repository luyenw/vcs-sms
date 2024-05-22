package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"vcs-sms/config"
	"vcs-sms/model/dto"
)

type MockReportService struct {
	mock.Mock
}

func (m *MockReportService) SendReport(startMils int64, endMils int64, to []string) error {
	args := m.Called(startMils, endMils, to)
	return args.Error(0)
}
func (m *MockReportService) PeriodicReport(interval time.Duration) {
	m.Called(interval)
}
func TestSendReport(t *testing.T) {
	t.Run("Failed to Bind JSON", func(t *testing.T) {
		mockReportService := new(MockReportService)
		reportController := NewReportController(mockReportService)
		router := config.GetTestGin()
		router.POST("/sendReport", reportController.SendReport)
		jsonValue := `{"StartDate": "01-01-2023", "EndDate": "31-01-2023", "Email": "`
		req, _ := http.NewRequest("POST", "/sendReport", bytes.NewBuffer([]byte(jsonValue)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Invalid Start Date Format", func(t *testing.T) {
		mockReportService := new(MockReportService)
		reportController := NewReportController(mockReportService)
		router := config.GetTestGin()
		router.POST("/sendReport", reportController.SendReport)
		requestBody := dto.ReportRequest{
			EndDate:   "2023-01-01",
			StartDate: "31-01-2023",
			Email:     "test@example.com",
		}
		jsonValue, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/sendReport", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid End Date Format", func(t *testing.T) {
		mockReportService := new(MockReportService)
		reportController := NewReportController(mockReportService)
		router := config.GetTestGin()
		router.POST("/sendReport", reportController.SendReport)
		requestBody := dto.ReportRequest{
			StartDate: "01-01-2023",
			EndDate:   "2023-01-31",
			Email:     "test@example.com",
		}
		jsonValue, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/sendReport", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Failed to Send Report", func(t *testing.T) {
		mockReportService := new(MockReportService)
		reportController := NewReportController(mockReportService)
		router := config.GetTestGin()
		router.POST("/sendReport", reportController.SendReport)
		requestBody := dto.ReportRequest{
			StartDate: "01-01-2023",
			EndDate:   "01-01-2025",
			Email:     "test@example.com",
		}
		jsonValue, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/sendReport", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		mockReportService.On("SendReport", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Failed to send report"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("Success", func(t *testing.T) {
		mockReportService := new(MockReportService)
		reportController := NewReportController(mockReportService)
		router := config.GetTestGin()
		router.POST("/sendReport", reportController.SendReport)
		requestBody := dto.ReportRequest{
			StartDate: "01-01-2023",
			EndDate:   "31-01-2023",
			Email:     "test@example.com",
		}
		jsonValue, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/sendReport", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		mockReportService.On("SendReport", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
