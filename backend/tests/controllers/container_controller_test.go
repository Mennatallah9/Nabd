package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"nabd/interfaces"
	"nabd/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDockerService struct {
	mock.Mock
}

func (m *MockDockerService) GetContainers() ([]models.ContainerInfo, error) {
	args := m.Called()
	return args.Get(0).([]models.ContainerInfo), args.Error(1)
}

func (m *MockDockerService) GetContainerLogs(containerName string, lines int) ([]string, error) {
	args := m.Called(containerName, lines)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockDockerService) RestartContainer(containerName string) error {
	args := m.Called(containerName)
	return args.Error(0)
}

func (m *MockDockerService) GetContainerMetrics() ([]models.ContainerMetric, error) {
	args := m.Called()
	return args.Get(0).([]models.ContainerMetric), args.Error(1)
}

type MockMetricsService struct {
	mock.Mock
}

func (m *MockMetricsService) GetLatestMetrics() ([]models.ContainerMetric, error) {
	args := m.Called()
	return args.Get(0).([]models.ContainerMetric), args.Error(1)
}

func (m *MockMetricsService) GetMetricsHistory(containerID string, hours int) ([]models.ContainerMetric, error) {
	args := m.Called(containerID, hours)
	return args.Get(0).([]models.ContainerMetric), args.Error(1)
}

func (m *MockMetricsService) CollectAndStoreMetrics() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMetricsService) GetActiveAlerts() ([]models.Alert, error) {
	args := m.Called()
	return args.Get(0).([]models.Alert), args.Error(1)
}

var _ interfaces.DockerServiceInterface = (*MockDockerService)(nil)
var _ interfaces.MetricsServiceInterface = (*MockMetricsService)(nil)

func createTestController(dockerService interfaces.DockerServiceInterface, metricsService interfaces.MetricsServiceInterface) *testControllerWrapper {
	return &testControllerWrapper{
		dockerService:  dockerService,
		metricsService: metricsService,
	}
}

type testControllerWrapper struct {
	dockerService  interfaces.DockerServiceInterface
	metricsService interfaces.MetricsServiceInterface
}

func (tc *testControllerWrapper) GetContainers(c *gin.Context) {
	containers, err := tc.dockerService.GetContainers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": containers})
}

func (tc *testControllerWrapper) GetMetrics(c *gin.Context) {
	metrics, err := tc.metricsService.GetLatestMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": metrics})
}

func (tc *testControllerWrapper) GetLogs(c *gin.Context) {
	containerName := c.Query("container")
	if containerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "container parameter is required"})
		return
	}
	
	lines := 100 // default
	if linesParam := c.Query("lines"); linesParam != "" {
		if linesParam == "50" {
			lines = 50
		}
	}

	logs, err := tc.dockerService.GetContainerLogs(containerName, lines)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

func (tc *testControllerWrapper) RestartContainer(c *gin.Context) {
	containerName := c.Param("name")
	
	err := tc.dockerService.RestartContainer(containerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container restarted successfully"})
}

func TestContainerController_GetContainers_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockDockerService := &MockDockerService{}
	mockMetricsService := &MockMetricsService{}
	
	expectedContainers := []models.ContainerInfo{
		{
			ID:      "container1",
			Name:    "test-container-1",
			Image:   "nginx:latest",
			Status:  "running",
			State:   "running",
			Created: time.Now(),
		},
		{
			ID:      "container2", 
			Name:    "test-container-2",
			Image:   "redis:alpine",
			Status:  "stopped",
			State:   "exited",
			Created: time.Now(),
		},
	}
	
	mockDockerService.On("GetContainers").Return(expectedContainers, nil)
	
	controller := createTestController(mockDockerService, mockMetricsService)
	
	router := gin.New()
	router.GET("/containers", controller.GetContainers)
	
	req := httptest.NewRequest("GET", "/containers", nil)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Contains(t, response, "data")
	mockDockerService.AssertExpectations(t)
}

func TestContainerController_GetContainers_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockDockerService := &MockDockerService{}
	mockMetricsService := &MockMetricsService{}
	
	mockDockerService.On("GetContainers").Return([]models.ContainerInfo{}, fmt.Errorf("docker service error"))
	
	controller := createTestController(mockDockerService, mockMetricsService)
	
	router := gin.New()
	router.GET("/containers", controller.GetContainers)
	
	req := httptest.NewRequest("GET", "/containers", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "docker service error")
	mockDockerService.AssertExpectations(t)
}

func TestContainerController_GetMetrics_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockDockerService := &MockDockerService{}
	mockMetricsService := &MockMetricsService{}

	expectedMetrics := []models.ContainerMetric{
		{
			ContainerID: "container1",
			Name:        "test-container-1",
			CPUPercent:  25.5,
			MemoryUsage: 1024000,
			MemoryLimit: 2048000,
			Status:      "running",
			Timestamp:   time.Now(),
		},
	}
	
	mockMetricsService.On("GetLatestMetrics").Return(expectedMetrics, nil)
	
	controller := createTestController(mockDockerService, mockMetricsService)

	router := gin.New()
	router.GET("/metrics", controller.GetMetrics)

	req := httptest.NewRequest("GET", "/metrics", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Contains(t, response, "data")
	mockMetricsService.AssertExpectations(t)
}

func TestContainerController_GetMetrics_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockDockerService := &MockDockerService{}
	mockMetricsService := &MockMetricsService{}

	mockMetricsService.On("GetLatestMetrics").Return([]models.ContainerMetric{}, fmt.Errorf("metrics service error"))
	
	controller := createTestController(mockDockerService, mockMetricsService)

	router := gin.New()
	router.GET("/metrics", controller.GetMetrics)

	req := httptest.NewRequest("GET", "/metrics", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "metrics service error")
	mockMetricsService.AssertExpectations(t)
}

func TestContainerController_GetLogs_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockDockerService := &MockDockerService{}
	mockMetricsService := &MockMetricsService{}

	expectedLogs := []string{
		"2024-01-01 12:00:00 [INFO] Application started",
		"2024-01-01 12:00:01 [INFO] Processing request",
		"2024-01-01 12:00:02 [INFO] Request completed",
	}
	
	mockDockerService.On("GetContainerLogs", "test-container", 100).Return(expectedLogs, nil)
	
	controller := createTestController(mockDockerService, mockMetricsService)

	router := gin.New()
	router.GET("/logs", controller.GetLogs)

	req := httptest.NewRequest("GET", "/logs?container=test-container", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Contains(t, response, "data")
	mockDockerService.AssertExpectations(t)
}

func TestContainerController_GetLogs_MissingContainer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockDockerService := &MockDockerService{}
	mockMetricsService := &MockMetricsService{}
	
	controller := createTestController(mockDockerService, mockMetricsService)

	router := gin.New()
	router.GET("/logs", controller.GetLogs)

	req := httptest.NewRequest("GET", "/logs", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "container parameter is required")
}

func TestContainerController_RestartContainer_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockDockerService := &MockDockerService{}
	mockMetricsService := &MockMetricsService{}
	
	mockDockerService.On("RestartContainer", "test-container").Return(nil)
	
	controller := createTestController(mockDockerService, mockMetricsService)

	router := gin.New()
	router.POST("/restart/:name", controller.RestartContainer)

	req := httptest.NewRequest("POST", "/restart/test-container", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Container restarted successfully")
	mockDockerService.AssertExpectations(t)
}

func TestContainerController_RestartContainer_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockDockerService := &MockDockerService{}
	mockMetricsService := &MockMetricsService{}
	
	mockDockerService.On("RestartContainer", "test-container").Return(fmt.Errorf("restart failed"))
	
	controller := createTestController(mockDockerService, mockMetricsService)

	router := gin.New()
	router.POST("/restart/:name", controller.RestartContainer)

	req := httptest.NewRequest("POST", "/restart/test-container", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "restart failed")
	mockDockerService.AssertExpectations(t)
}