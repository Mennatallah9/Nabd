package services

import (
	"log"
	"nabd/models"
	"time"
)

type AutoHealService struct {
	dockerService  *DockerService
	metricsService *MetricsService
	config         *models.Config
}

// NewAutoHealService creates a new auto-heal service
func NewAutoHealService(dockerService *DockerService, metricsService *MetricsService, config *models.Config) *AutoHealService {
	return &AutoHealService{
		dockerService:  dockerService,
		metricsService: metricsService,
		config:         config,
	}
}

// StartAutoHealing starts the auto-healing process
func (ahs *AutoHealService) StartAutoHealing() {
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	go func() {
		for range ticker.C {
			ahs.PerformAutoHealing()
		}
	}()
	log.Println("Auto-healing service started")
}

// PerformAutoHealing checks for unhealthy containers and heals them
func (ahs *AutoHealService) PerformAutoHealing() {
	events := ahs.dockerService.CheckUnhealthyContainers()
	
	for _, event := range events {
		if err := ahs.storeAutoHealEvent(event); err != nil {
			log.Printf("Error storing auto-heal event: %v", err)
		}
	}
	
	if len(events) > 0 {
		log.Printf("Auto-healing completed: %d actions performed", len(events))
	}
}

// storeAutoHealEvent stores an auto-heal event in the database
func (ahs *AutoHealService) storeAutoHealEvent(event models.AutoHealEvent) error {
	query := `INSERT INTO autoheal_events 
		(container_id, name, action, reason, success, timestamp)
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err := models.DB.Exec(query,
		event.ContainerID,
		event.Name,
		event.Action,
		event.Reason,
		event.Success,
		event.Timestamp,
	)

	return err
}

// GetAutoHealHistory returns recent auto-heal events
func (ahs *AutoHealService) GetAutoHealHistory(limit int) ([]models.AutoHealEvent, error) {
	query := `SELECT id, container_id, name, action, reason, success, timestamp
		FROM autoheal_events 
		ORDER BY timestamp DESC 
		LIMIT ?`

	rows, err := models.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.AutoHealEvent
	for rows.Next() {
		var event models.AutoHealEvent
		err := rows.Scan(
			&event.ID,
			&event.ContainerID,
			&event.Name,
			&event.Action,
			&event.Reason,
			&event.Success,
			&event.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}