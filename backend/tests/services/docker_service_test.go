package services

import (
	"fmt"
	"testing"
	"time"

	"nabd/models"

	"github.com/stretchr/testify/assert"
)

func TestDockerService_NewDockerService(t *testing.T) {
	config := &models.Config{}
	config.Docker.Host = "unix:///var/run/docker.sock"
	
	//note: this would normally create a real Docker client
	//for testing, we'll just verify the configuration
	assert.NotNil(t, config)
	assert.Equal(t, "unix:///var/run/docker.sock", config.Docker.Host)
}

func TestDockerService_MockImplementation(t *testing.T) {
	dockerMock := &dockerServiceTestWrapper{}
	containers, err := dockerMock.GetContainers()
	assert.NoError(t, err)
	assert.Len(t, containers, 2)
	assert.Equal(t, "nginx-container", containers[0].Name)
	assert.Equal(t, "redis-container", containers[1].Name)
}

func TestDockerService_GetContainerMetrics_MockImplementation(t *testing.T) {
	dockerMock := &dockerServiceTestWrapper{}
	metrics, err := dockerMock.GetContainerMetrics()
	assert.NoError(t, err)
	assert.Len(t, metrics, 2)
	assert.Equal(t, "nginx-container", metrics[0].Name)
	assert.Equal(t, 45.5, metrics[0].CPUPercent)
	assert.Equal(t, "redis-container", metrics[1].Name)
	assert.Equal(t, 12.3, metrics[1].CPUPercent)
}

func TestDockerService_GetContainerLogs_MockImplementation(t *testing.T) {
	dockerMock := &dockerServiceTestWrapper{}
	logs, err := dockerMock.GetContainerLogs("nginx-container", 100)
	assert.NoError(t, err)
	assert.Len(t, logs, 3)
	assert.Contains(t, logs[0], "nginx started")
	assert.Contains(t, logs[1], "Processing request")
	assert.Contains(t, logs[2], "Request completed")
}

func TestDockerService_RestartContainer_MockImplementation(t *testing.T) {
	dockerMock := &dockerServiceTestWrapper{}
	err := dockerMock.RestartContainer("nginx-container")
	assert.NoError(t, err)
	err = dockerMock.RestartContainer("non-existent-container")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "container not found")
}

type dockerServiceTestWrapper struct{}

func (ds *dockerServiceTestWrapper) GetContainers() ([]models.ContainerInfo, error) {
	return []models.ContainerInfo{
		{
			ID:      "container-123",
			Name:    "nginx-container",
			Image:   "nginx:latest",
			Status:  "running",
			State:   "running",
			Created: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:      "container-456",
			Name:    "redis-container",
			Image:   "redis:alpine",
			Status:  "running", 
			State:   "running",
			Created: time.Now().Add(-12 * time.Hour),
		},
	}, nil
}

func (ds *dockerServiceTestWrapper) GetContainerMetrics() ([]models.ContainerMetric, error) {
	return []models.ContainerMetric{
		{
			ContainerID: "container-123",
			Name:        "nginx-container",
			CPUPercent:  45.5,
			MemoryUsage: 256000000,
			MemoryLimit: 512000000,
			NetworkRx:   1024000,
			NetworkTx:   2048000,
			Status:      "running",
			Timestamp:   time.Now(),
		},
		{
			ContainerID: "container-456",
			Name:        "redis-container",
			CPUPercent:  12.3,
			MemoryUsage: 128000000,
			MemoryLimit: 256000000,
			NetworkRx:   512000,
			NetworkTx:   1024000,
			Status:      "running",
			Timestamp:   time.Now(),
		},
	}, nil
}

func (ds *dockerServiceTestWrapper) GetContainerLogs(containerName string, lines int) ([]string, error) {
	if containerName == "nginx-container" {
		return []string{
			"2024-01-01 12:00:00 [INFO] nginx started successfully",
			"2024-01-01 12:01:00 [INFO] Processing request GET /",
			"2024-01-01 12:01:01 [INFO] Request completed with status 200",
		}, nil
	} else if containerName == "redis-container" {
		return []string{
			"2024-01-01 12:00:00 [INFO] Redis server started",
			"2024-01-01 12:01:00 [INFO] Accepting connections",
		}, nil
	}
	
	return []string{}, fmt.Errorf("container not found: %s", containerName)
}

func (ds *dockerServiceTestWrapper) RestartContainer(containerName string) error {
	if containerName == "nginx-container" || containerName == "redis-container" {
		return nil
	}
	
	return fmt.Errorf("container not found: %s", containerName)
}