package models_test

import (
	"encoding/json"
	"testing"
	"time"

	"nabd/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContainerMetric_Struct(t *testing.T) {
	metric := models.ContainerMetric{
		ID:          1,
		ContainerID: "container123",
		Name:        "test-container",
		CPUPercent:  75.5,
		MemoryUsage: 1024000,
		MemoryLimit: 2048000,
		NetworkRx:   5000,
		NetworkTx:   3000,
		Status:      "running",
		Timestamp:   time.Now(),
	}

	assert.Equal(t, 1, metric.ID)
	assert.Equal(t, "container123", metric.ContainerID)
	assert.Equal(t, "test-container", metric.Name)
	assert.Equal(t, 75.5, metric.CPUPercent)
	assert.Equal(t, int64(1024000), metric.MemoryUsage)
	assert.Equal(t, int64(2048000), metric.MemoryLimit)
	assert.Equal(t, int64(5000), metric.NetworkRx)
	assert.Equal(t, int64(3000), metric.NetworkTx)
	assert.Equal(t, "running", metric.Status)
}

func TestAutoHealEvent_MarshalJSON(t *testing.T) {
	timestamp := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	
	event := models.AutoHealEvent{
		ID:          1,
		ContainerID: "container123",
		Name:        "test-container",
		Action:      "restart",
		Reason:      "high CPU usage",
		Success:     true,
		Timestamp:   timestamp,
	}

	jsonData, err := json.Marshal(event)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	assert.Equal(t, float64(1), result["id"])
	assert.Equal(t, "container123", result["container_id"])
	assert.Equal(t, "test-container", result["name"])
	assert.Equal(t, "restart", result["action"])
	assert.Equal(t, "high CPU usage", result["reason"])
	assert.Equal(t, true, result["success"])
	assert.Equal(t, "2024-01-01T12:00:00Z", result["timestamp"])
}

func TestAlert_Struct(t *testing.T) {
	alert := models.Alert{
		ID:          1,
		ContainerID: "container123",
		Name:        "test-container",
		Type:        "cpu",
		Message:     "High CPU usage detected",
		Severity:    "warning",
		Active:      true,
		Timestamp:   time.Now(),
	}

	assert.Equal(t, 1, alert.ID)
	assert.Equal(t, "container123", alert.ContainerID)
	assert.Equal(t, "test-container", alert.Name)
	assert.Equal(t, "cpu", alert.Type)
	assert.Equal(t, "High CPU usage detected", alert.Message)
	assert.Equal(t, "warning", alert.Severity)
	assert.True(t, alert.Active)
}

func TestContainerInfo_Struct(t *testing.T) {
	info := models.ContainerInfo{
		ID:      "container123",
		Name:    "test-container",
		Image:   "nginx:latest",
		Status:  "running",
		State:   "running",
		Created: time.Now(),
	}

	assert.Equal(t, "container123", info.ID)
	assert.Equal(t, "test-container", info.Name)
	assert.Equal(t, "nginx:latest", info.Image)
	assert.Equal(t, "running", info.Status)
	assert.Equal(t, "running", info.State)
}

func TestConfig_Struct(t *testing.T) {
	config := models.Config{}
	config.Database.Path = "/path/to/db"
	config.Docker.Host = "unix:///var/run/docker.sock"
	config.Auth.AdminToken = "admin-token"
	config.Alerts.CPUThreshold = 90.0
	config.Alerts.MemoryThreshold = 85.0
	config.Alerts.RestartLimit = 3
	config.AutoHeal.Enabled = true
	config.AutoHeal.Interval = 60

	assert.Equal(t, "/path/to/db", config.Database.Path)
	assert.Equal(t, "unix:///var/run/docker.sock", config.Docker.Host)
	assert.Equal(t, "admin-token", config.Auth.AdminToken)
	assert.Equal(t, 90.0, config.Alerts.CPUThreshold)
	assert.Equal(t, 85.0, config.Alerts.MemoryThreshold)
	assert.Equal(t, 3, config.Alerts.RestartLimit)
	assert.True(t, config.AutoHeal.Enabled)
	assert.Equal(t, 60, config.AutoHeal.Interval)
}