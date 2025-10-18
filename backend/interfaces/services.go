package interfaces

import "nabd/models"

type DockerServiceInterface interface {
	GetContainers() ([]models.ContainerInfo, error)
	GetContainerLogs(containerName string, lines int) ([]string, error)
	RestartContainer(containerName string) error
	GetContainerMetrics() ([]models.ContainerMetric, error)
}

type MetricsServiceInterface interface {
	GetLatestMetrics() ([]models.ContainerMetric, error)
	GetMetricsHistory(containerID string, hours int) ([]models.ContainerMetric, error)
	CollectAndStoreMetrics() error
	GetActiveAlerts() ([]models.Alert, error)
}