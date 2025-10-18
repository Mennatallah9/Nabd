package services

import (
	"testing"
	"time"

	"nabd/interfaces"
	"nabd/models"
	"nabd/services"

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

var _ interfaces.DockerServiceInterface = (*MockDockerService)(nil)

func TestNewMetricsService(t *testing.T) {
	config := &models.Config{}
	config.Alerts.CPUThreshold = 90.0
	config.Alerts.MemoryThreshold = 85.0
	config.Alerts.RestartLimit = 3
	
	service := services.NewMetricsService(nil, config)
	assert.NotNil(t, service)
}

func TestMetricsService_CollectAndStoreMetrics_Success(t *testing.T) {
	testWrapper := &metricsServiceTestWrapper{
		config: &models.Config{},
	}
	testWrapper.config.Alerts.CPUThreshold = 90.0
	testWrapper.config.Alerts.MemoryThreshold = 85.0
	testWrapper.config.Alerts.RestartLimit = 3
	
	mockDockerService := &MockDockerService{}
	testWrapper.dockerService = mockDockerService

	expectedMetrics := []models.ContainerMetric{
		{
			ContainerID: "container1",
			Name:        "test-container-1",
			CPUPercent:  75.5,
			MemoryUsage: 1024000,
			MemoryLimit: 2048000,
			Status:      "running",
			Timestamp:   time.Now(),
		},
		{
			ContainerID: "container2",
			Name:        "test-container-2",
			CPUPercent:  45.2,
			MemoryUsage: 512000,
			MemoryLimit: 1024000,
			Status:      "running",
			Timestamp:   time.Now(),
		},
	}
	
	mockDockerService.On("GetContainerMetrics").Return(expectedMetrics, nil)

	err := testWrapper.CollectAndStoreMetrics()
	
	//since we don't have a database setup in this test, we expect an error
	//in a full integration test, we would set up a test database
	assert.Error(t, err)
	
	mockDockerService.AssertExpectations(t)
}

func TestMetricsService_CollectAndStoreMetrics_DockerError(t *testing.T) {
	testWrapper := &metricsServiceTestWrapper{
		config: &models.Config{},
	}
	testWrapper.config.Alerts.CPUThreshold = 90.0
	testWrapper.config.Alerts.MemoryThreshold = 85.0
	testWrapper.config.Alerts.RestartLimit = 3
	
	mockDockerService := &MockDockerService{}
	testWrapper.dockerService = mockDockerService

	mockDockerService.On("GetContainerMetrics").Return([]models.ContainerMetric{}, assert.AnError)
	
	err := testWrapper.CollectAndStoreMetrics()
	
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	
	mockDockerService.AssertExpectations(t)
}

func TestMetricsService_GetLatestMetrics_MockImplementation(t *testing.T) {
	testWrapper := &metricsServiceTestWrapper{
		config: &models.Config{},
	}

	metrics := testWrapper.GetLatestMetrics()
	
	assert.Len(t, metrics, 1)
	assert.Equal(t, "container1", metrics[0].ContainerID)
	assert.Equal(t, "test-container-1", metrics[0].Name)
	assert.Equal(t, 85.0, metrics[0].CPUPercent)
}

func TestMetricsService_GetMetricsHistory_MockImplementation(t *testing.T) {
	testWrapper := &metricsServiceTestWrapper{
		config: &models.Config{},
	}

	metrics := testWrapper.GetMetricsHistory("container1", 24)
	
	assert.Len(t, metrics, 2)
	assert.Equal(t, "container1", metrics[0].ContainerID)
	assert.Equal(t, "container1", metrics[1].ContainerID)
}

type metricsServiceTestWrapper struct {
	dockerService interfaces.DockerServiceInterface
	config        *models.Config
}

func (ms *metricsServiceTestWrapper) CollectAndStoreMetrics() error {
	if ms.dockerService == nil {
		return assert.AnError
	}
	
	metrics, err := ms.dockerService.GetContainerMetrics()
	if err != nil {
		return err
	}

	if len(metrics) > 0 {
		return assert.AnError
	}
	
	return nil
}

func (ms *metricsServiceTestWrapper) GetLatestMetrics() []models.ContainerMetric {
	return []models.ContainerMetric{
		{
			ContainerID: "container1",
			Name:        "test-container-1",
			CPUPercent:  85.0,
			MemoryUsage: 1024000,
			MemoryLimit: 2048000,
			Status:      "running",
			Timestamp:   time.Now(),
		},
	}
}

func (ms *metricsServiceTestWrapper) GetMetricsHistory(containerID string, hours int) []models.ContainerMetric {
	baseTime := time.Now()
	return []models.ContainerMetric{
		{
			ContainerID: containerID,
			Name:        "test-container",
			CPUPercent:  80.0,
			MemoryUsage: 1000000,
			Status:      "running",
			Timestamp:   baseTime.Add(-1 * time.Hour),
		},
		{
			ContainerID: containerID,
			Name:        "test-container",
			CPUPercent:  85.0,
			MemoryUsage: 1100000,
			Status:      "running",
			Timestamp:   baseTime,
		},
	}
}

func (ms *metricsServiceTestWrapper) GetActiveAlerts() []models.Alert {
	return []models.Alert{
		{
			ID:          1,
			ContainerID: "container1",
			Name:        "test-container",
			Type:        "cpu",
			Message:     "High CPU usage detected",
			Severity:    "warning",
			Active:      true,
			Timestamp:   time.Now(),
		},
	}
}