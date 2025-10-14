package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type ContainerMetric struct {
	ID          int       `json:"id" db:"id"`
	ContainerID string    `json:"container_id" db:"container_id"`
	Name        string    `json:"name" db:"name"`
	CPUPercent  float64   `json:"cpu_percent" db:"cpu_percent"`
	MemoryUsage int64     `json:"memory_usage" db:"memory_usage"`
	MemoryLimit int64     `json:"memory_limit" db:"memory_limit"`
	NetworkRx   int64     `json:"network_rx" db:"network_rx"`
	NetworkTx   int64     `json:"network_tx" db:"network_tx"`
	Status      string    `json:"status" db:"status"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
}

type AutoHealEvent struct {
	ID          int       `json:"id" db:"id"`
	ContainerID string    `json:"container_id" db:"container_id"`
	Name        string    `json:"name" db:"name"`
	Action      string    `json:"action" db:"action"`
	Reason      string    `json:"reason" db:"reason"`
	Success     bool      `json:"success" db:"success"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
}

// custom marshaling for AutoHealEvent to ensure proper timestamp format
func (ahe AutoHealEvent) MarshalJSON() ([]byte, error) {
	type Alias AutoHealEvent
	return json.Marshal(&struct {
		Timestamp string `json:"timestamp"`
		*Alias
	}{
		Timestamp: ahe.Timestamp.Format(time.RFC3339),
		Alias:     (*Alias)(&ahe),
	})
}

type Alert struct {
	ID          int       `json:"id" db:"id"`
	ContainerID string    `json:"container_id" db:"container_id"`
	Name        string    `json:"name" db:"name"`
	Type        string    `json:"type" db:"type"`
	Message     string    `json:"message" db:"message"`
	Severity    string    `json:"severity" db:"severity"`
	Active      bool      `json:"active" db:"active"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
}

type ContainerInfo struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Image   string    `json:"image"`
	Status  string    `json:"status"`
	State   string    `json:"state"`
	Created time.Time `json:"created"`
}

type Config struct {
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
	Docker struct {
		Host string `yaml:"host"`
	} `yaml:"docker"`
	Auth struct {
		AdminToken string `yaml:"admin_token"`
	} `yaml:"auth"`
	AutoHeal struct {
		Enabled           bool     `yaml:"enabled"`
		Interval          int      `yaml:"interval"`
		ExcludeContainers []string `yaml:"exclude_containers"`
	} `yaml:"autoheal"`
	Alerts struct {
		CPUThreshold    float64 `yaml:"cpu_threshold"`
		MemoryThreshold float64 `yaml:"memory_threshold"`
		RestartLimit    int     `yaml:"restart_limit"`
	} `yaml:"alerts"`
}

// Database connection
var DB *sql.DB