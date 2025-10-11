package services

import (
	"database/sql"
	"log"
	"nabd/models"
	"time"
)

type MetricsService struct {
	dockerService *DockerService
	config        *models.Config
}

// NewMetricsService creates a new metrics service
func NewMetricsService(dockerService *DockerService, config *models.Config) *MetricsService {
	return &MetricsService{
		dockerService: dockerService,
		config:        config,
	}
}

// CollectAndStoreMetrics collects metrics from Docker and stores them in the database
func (ms *MetricsService) CollectAndStoreMetrics() error {
	metrics, err := ms.dockerService.GetContainerMetrics()
	if err != nil {
		return err
	}

	for _, metric := range metrics {
		if err := ms.storeMetric(metric); err != nil {
			log.Printf("Error storing metric for container %s: %v", metric.Name, err)
		}

		// Check for alerts
		if err := ms.checkAlerts(metric); err != nil {
			log.Printf("Error checking alerts for container %s: %v", metric.Name, err)
		}
	}

	return nil
}

// storeMetric stores a single metric in the database
func (ms *MetricsService) storeMetric(metric models.ContainerMetric) error {
	query := `INSERT INTO container_metrics 
		(container_id, name, cpu_percent, memory_usage, memory_limit, network_rx, network_tx, status, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := models.DB.Exec(query,
		metric.ContainerID,
		metric.Name,
		metric.CPUPercent,
		metric.MemoryUsage,
		metric.MemoryLimit,
		metric.NetworkRx,
		metric.NetworkTx,
		metric.Status,
		metric.Timestamp,
	)

	return err
}

// GetLatestMetrics returns the latest metrics for all containers
func (ms *MetricsService) GetLatestMetrics() ([]models.ContainerMetric, error) {
	query := `SELECT DISTINCT 
		container_id, name, cpu_percent, memory_usage, memory_limit, 
		network_rx, network_tx, status, timestamp
		FROM container_metrics cm1
		WHERE timestamp = (
			SELECT MAX(timestamp) 
			FROM container_metrics cm2 
			WHERE cm2.container_id = cm1.container_id
		)
		ORDER BY name`

	rows, err := models.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []models.ContainerMetric
	for rows.Next() {
		var metric models.ContainerMetric
		err := rows.Scan(
			&metric.ContainerID,
			&metric.Name,
			&metric.CPUPercent,
			&metric.MemoryUsage,
			&metric.MemoryLimit,
			&metric.NetworkRx,
			&metric.NetworkTx,
			&metric.Status,
			&metric.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// GetMetricsHistory returns historical metrics for a container
func (ms *MetricsService) GetMetricsHistory(containerID string, hours int) ([]models.ContainerMetric, error) {
	query := `SELECT container_id, name, cpu_percent, memory_usage, memory_limit,
		network_rx, network_tx, status, timestamp
		FROM container_metrics 
		WHERE container_id = ? AND timestamp > datetime('now', '-' || ? || ' hours')
		ORDER BY timestamp DESC`

	rows, err := models.DB.Query(query, containerID, hours)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []models.ContainerMetric
	for rows.Next() {
		var metric models.ContainerMetric
		err := rows.Scan(
			&metric.ContainerID,
			&metric.Name,
			&metric.CPUPercent,
			&metric.MemoryUsage,
			&metric.MemoryLimit,
			&metric.NetworkRx,
			&metric.NetworkTx,
			&metric.Status,
			&metric.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// checkAlerts checks if metrics trigger any alerts
func (ms *MetricsService) checkAlerts(metric models.ContainerMetric) error {
	// Check CPU alert
	if metric.CPUPercent > ms.config.Alerts.CPUThreshold {
		alert := models.Alert{
			ContainerID: metric.ContainerID,
			Name:        metric.Name,
			Type:        "high_cpu",
			Message:     "High CPU usage detected",
			Severity:    "warning",
			Active:      true,
			Timestamp:   time.Now(),
		}
		if err := ms.storeAlert(alert); err != nil {
			return err
		}
	}

	// Check memory alert
	if metric.MemoryLimit > 0 {
		memoryPercent := float64(metric.MemoryUsage) / float64(metric.MemoryLimit) * 100
		if memoryPercent > ms.config.Alerts.MemoryThreshold {
			alert := models.Alert{
				ContainerID: metric.ContainerID,
				Name:        metric.Name,
				Type:        "high_memory",
				Message:     "High memory usage detected",
				Severity:    "warning",
				Active:      true,
				Timestamp:   time.Now(),
			}
			if err := ms.storeAlert(alert); err != nil {
				return err
			}
		}
	}

	return nil
}

// storeAlert stores an alert in the database
func (ms *MetricsService) storeAlert(alert models.Alert) error {
	// Check if similar alert already exists and is active
	var count int
	checkQuery := `SELECT COUNT(*) FROM alerts 
		WHERE container_id = ? AND type = ? AND active = 1 
		AND timestamp > datetime('now', '-1 hour')`
	
	err := models.DB.QueryRow(checkQuery, alert.ContainerID, alert.Type).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// Don't create duplicate alerts
	if count > 0 {
		return nil
	}

	query := `INSERT INTO alerts 
		(container_id, name, type, message, severity, active, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = models.DB.Exec(query,
		alert.ContainerID,
		alert.Name,
		alert.Type,
		alert.Message,
		alert.Severity,
		alert.Active,
		alert.Timestamp,
	)

	return err
}

// GetActiveAlerts returns all active alerts
func (ms *MetricsService) GetActiveAlerts() ([]models.Alert, error) {
	query := `SELECT id, container_id, name, type, message, severity, active, timestamp
		FROM alerts 
		WHERE active = 1 
		ORDER BY timestamp DESC`

	rows, err := models.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []models.Alert
	for rows.Next() {
		var alert models.Alert
		err := rows.Scan(
			&alert.ID,
			&alert.ContainerID,
			&alert.Name,
			&alert.Type,
			&alert.Message,
			&alert.Severity,
			&alert.Active,
			&alert.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}