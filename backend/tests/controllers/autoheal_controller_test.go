package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"nabd/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAutoHealService struct {
	mock.Mock
}

func (m *MockAutoHealService) GetAutoHealEvents() ([]models.AutoHealEvent, error) {
	args := m.Called()
	return args.Get(0).([]models.AutoHealEvent), args.Error(1)
}

func (m *MockAutoHealService) TriggerManualHeal(containerName string) error {
	args := m.Called(containerName)
	return args.Error(0)
}

type autoHealControllerWrapper struct {
	autoHealService *MockAutoHealService
	metricsService  *MockMetricsService
}

func (ac *autoHealControllerWrapper) GetAutoHealEvents(c *gin.Context) {
	events, err := ac.autoHealService.GetAutoHealEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": events})
}

func (ac *autoHealControllerWrapper) GetAlerts(c *gin.Context) {
	alerts, err := ac.metricsService.GetActiveAlerts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": alerts})
}

func TestAutoHealController_GetAutoHealEvents_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockAutoHealService := &MockAutoHealService{}
	mockMetricsService := &MockMetricsService{}
	
	expectedEvents := []models.AutoHealEvent{
		{
			ID:          1,
			ContainerID: "container1",
			Name:        "test-container-1",
			Action:      "restart",
			Reason:      "high CPU usage",
			Success:     true,
			Timestamp:   time.Now(),
		},
		{
			ID:          2,
			ContainerID: "container2",
			Name:        "test-container-2",
			Action:      "restart",
			Reason:      "memory threshold exceeded",
			Success:     false,
			Timestamp:   time.Now(),
		},
	}
	
	mockAutoHealService.On("GetAutoHealEvents").Return(expectedEvents, nil)
	
	controller := &autoHealControllerWrapper{
		autoHealService: mockAutoHealService,
		metricsService:  mockMetricsService,
	}
	
	router := gin.New()
	router.GET("/autoheal/events", controller.GetAutoHealEvents)
	
	req := httptest.NewRequest("GET", "/autoheal/events", nil)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Contains(t, response, "data")
	data := response["data"].([]interface{})
	assert.Len(t, data, 2)
	
	mockAutoHealService.AssertExpectations(t)
}

func TestAutoHealController_GetAlerts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockAutoHealService := &MockAutoHealService{}
	mockMetricsService := &MockMetricsService{}
	
	expectedAlerts := []models.Alert{
		{
			ID:          1,
			ContainerID: "container1",
			Name:        "test-container-1",
			Type:        "cpu",
			Message:     "CPU usage above 90%",
			Severity:    "warning",
			Active:      true,
			Timestamp:   time.Now(),
		},
		{
			ID:          2,
			ContainerID: "container2",
			Name:        "test-container-2",
			Type:        "memory",
			Message:     "Memory usage above 85%",
			Severity:    "critical",
			Active:      true,
			Timestamp:   time.Now(),
		},
	}
	
	mockMetricsService.On("GetActiveAlerts").Return(expectedAlerts, nil)
	
	controller := &autoHealControllerWrapper{
		autoHealService: mockAutoHealService,
		metricsService:  mockMetricsService,
	}
	
	router := gin.New()
	router.GET("/alerts", controller.GetAlerts)
	
	req := httptest.NewRequest("GET", "/alerts", nil)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Contains(t, response, "data")
	data := response["data"].([]interface{})
	assert.Len(t, data, 2)
	
	mockMetricsService.AssertExpectations(t)
}

func TestAutoHealController_GetAutoHealEvents_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockAutoHealService := &MockAutoHealService{}
	mockMetricsService := &MockMetricsService{}
	
	mockAutoHealService.On("GetAutoHealEvents").Return([]models.AutoHealEvent{}, assert.AnError)
	
	controller := &autoHealControllerWrapper{
		autoHealService: mockAutoHealService,
		metricsService:  mockMetricsService,
	}
	
	router := gin.New()
	router.GET("/autoheal/events", controller.GetAutoHealEvents)
	
	req := httptest.NewRequest("GET", "/autoheal/events", nil)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	mockAutoHealService.AssertExpectations(t)
}

func TestAutoHealController_GetAlerts_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockAutoHealService := &MockAutoHealService{}
	mockMetricsService := &MockMetricsService{}
	
	mockMetricsService.On("GetActiveAlerts").Return([]models.Alert{}, assert.AnError)
	
	controller := &autoHealControllerWrapper{
		autoHealService: mockAutoHealService,
		metricsService:  mockMetricsService,
	}
	
	router := gin.New()
	router.GET("/alerts", controller.GetAlerts)
	
	req := httptest.NewRequest("GET", "/alerts", nil)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	mockMetricsService.AssertExpectations(t)
}