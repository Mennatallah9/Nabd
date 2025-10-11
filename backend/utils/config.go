package utils

import (
	"io/ioutil"
	"nabd/models"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads configuration from YAML file or environment variables
func LoadConfig() (*models.Config, error) {
	config := &models.Config{}

	// Set defaults
	config.Database.Path = "./nabd.db"
	config.Docker.Host = "unix:///var/run/docker.sock"
	config.Auth.AdminToken = "nabd-admin-token"
	config.Alerts.CPUThreshold = 90.0
	config.Alerts.MemoryThreshold = 90.0
	config.Alerts.RestartLimit = 3

	// Try to load from config file
	if _, err := os.Stat("config.yaml"); err == nil {
		data, err := ioutil.ReadFile("config.yaml")
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, err
		}
	}

	// Override with environment variables if present
	if token := os.Getenv("NABD_ADMIN_TOKEN"); token != "" {
		config.Auth.AdminToken = token
	}
	if dbPath := os.Getenv("NABD_DB_PATH"); dbPath != "" {
		config.Database.Path = dbPath
	}
	if dockerHost := os.Getenv("DOCKER_HOST"); dockerHost != "" {
		config.Docker.Host = dockerHost
	}

	return config, nil
}