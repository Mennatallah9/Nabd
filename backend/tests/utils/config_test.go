package utils

import (
	"os"
	"testing"

	"nabd/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_Defaults(t *testing.T) {
	originalWd, _ := os.Getwd()
	tempDir := t.TempDir()
	err := os.Chdir(tempDir)
	require.NoError(t, err)
	defer func() {
		os.Chdir(originalWd)
	}()

	config, err := utils.LoadConfig()
	
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "./nabd.db", config.Database.Path)
	assert.Equal(t, "unix:///var/run/docker.sock", config.Docker.Host)
	assert.Equal(t, "nabd-admin-token", config.Auth.AdminToken)
	assert.Equal(t, 90.0, config.Alerts.CPUThreshold)
	assert.Equal(t, 90.0, config.Alerts.MemoryThreshold)
	assert.Equal(t, 3, config.Alerts.RestartLimit)
}

func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
	originalAdminToken := os.Getenv("NABD_ADMIN_TOKEN")
	originalDBPath := os.Getenv("NABD_DB_PATH")
	os.Setenv("NABD_ADMIN_TOKEN", "test-token")
	os.Setenv("NABD_DB_PATH", "/test/db/path")

	defer func() {
		if originalAdminToken == "" {
			os.Unsetenv("NABD_ADMIN_TOKEN")
		} else {
			os.Setenv("NABD_ADMIN_TOKEN", originalAdminToken)
		}
		if originalDBPath == "" {
			os.Unsetenv("NABD_DB_PATH")
		} else {
			os.Setenv("NABD_DB_PATH", originalDBPath)
		}
	}()

	originalWd, _ := os.Getwd()
	tempDir := t.TempDir()
	err := os.Chdir(tempDir)
	require.NoError(t, err)
	defer func() {
		os.Chdir(originalWd)
	}()

	config, err := utils.LoadConfig()
	
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "test-token", config.Auth.AdminToken)
	assert.Equal(t, "/test/db/path", config.Database.Path)
}

func TestLoadConfig_WithConfigFile(t *testing.T) {
	originalWd, _ := os.Getwd()
	tempDir := t.TempDir()
	err := os.Chdir(tempDir)
	require.NoError(t, err)
	defer func() {
		os.Chdir(originalWd)
	}()

	configContent := `
database:
  path: "/custom/db/path"
auth:
  admin_token: "custom-token"
alerts:
  cpu_threshold: 80.0
  memory_threshold: 85.0
  restart_limit: 5
`
	err = os.WriteFile("config.yaml", []byte(configContent), 0644)
	require.NoError(t, err)

	config, err := utils.LoadConfig()
	
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "/custom/db/path", config.Database.Path)
	assert.Equal(t, "custom-token", config.Auth.AdminToken)
	assert.Equal(t, 80.0, config.Alerts.CPUThreshold)
	assert.Equal(t, 85.0, config.Alerts.MemoryThreshold)
	assert.Equal(t, 5, config.Alerts.RestartLimit)
}