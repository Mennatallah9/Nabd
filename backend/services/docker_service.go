package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"nabd/models"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerService struct {
	client *client.Client
	config *models.Config
}

// NewDockerService creates a new Docker service instance
func NewDockerService(config *models.Config) (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &DockerService{
		client: cli,
		config: config,
	}, nil
}

// GetContainers returns a list of all containers (excluding those in the exclusion list)
func (ds *DockerService) GetContainers() ([]models.ContainerInfo, error) {
	containers, err := ds.client.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var result []models.ContainerInfo
	for _, container := range containers {
		name := strings.TrimPrefix(container.Names[0], "/")
		
		// Check if container is in exclusion list
		excluded := false
		for _, excludedName := range ds.config.AutoHeal.ExcludeContainers {
			if name == excludedName {
				excluded = true
				break
			}
		}
		
		if excluded {
			continue
		}
		
		result = append(result, models.ContainerInfo{
			ID:      container.ID[:12],
			Name:    name,
			Image:   container.Image,
			Status:  container.Status,
			State:   container.State,
			Created: time.Unix(container.Created, 0),
		})
	}

	return result, nil
}

// GetContainerMetrics collects metrics for all running containers (excluding those in the exclusion list)
func (ds *DockerService) GetContainerMetrics() ([]models.ContainerMetric, error) {
	containers, err := ds.client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	var metrics []models.ContainerMetric
	for _, container := range containers {
		name := strings.TrimPrefix(container.Names[0], "/")
		
		// Check if container is in exclusion list
		excluded := false
		for _, excludedName := range ds.config.AutoHeal.ExcludeContainers {
			if name == excludedName {
				excluded = true
				break
			}
		}
		
		if excluded {
			continue
		}
		
		metric, err := ds.getContainerMetric(container)
		if err != nil {
			log.Printf("Error getting metrics for container %s: %v", container.ID[:12], err)
			continue
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// getContainerMetric gets metrics for a single container
func (ds *DockerService) getContainerMetric(container types.Container) (models.ContainerMetric, error) {
	name := strings.TrimPrefix(container.Names[0], "/")
	
	// Get container stats
	stats, err := ds.client.ContainerStats(context.Background(), container.ID, false)
	if err != nil {
		return models.ContainerMetric{}, err
	}
	defer stats.Body.Close()

	var statsData types.StatsJSON
	if err := json.NewDecoder(stats.Body).Decode(&statsData); err != nil {
		return models.ContainerMetric{}, err
	}

	// Calculate CPU percentage
	cpuPercent := calculateCPUPercent(&statsData)

	// Get memory usage
	memoryUsage := int64(statsData.MemoryStats.Usage)
	memoryLimit := int64(statsData.MemoryStats.Limit)

	// Get network stats
	var networkRx, networkTx int64
	for _, network := range statsData.Networks {
		networkRx += int64(network.RxBytes)
		networkTx += int64(network.TxBytes)
	}

	return models.ContainerMetric{
		ContainerID: container.ID[:12],
		Name:        name,
		CPUPercent:  cpuPercent,
		MemoryUsage: memoryUsage,
		MemoryLimit: memoryLimit,
		NetworkRx:   networkRx,
		NetworkTx:   networkTx,
		Status:      container.Status,
		Timestamp:   time.Now(),
	}, nil
}

// calculateCPUPercent calculates CPU usage percentage
func calculateCPUPercent(stats *types.StatsJSON) float64 {
	if stats.PreCPUStats.CPUUsage.TotalUsage == 0 {
		return 0.0
	}

	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		return (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return 0.0
}

// GetContainerLogs gets recent logs for a container
func (ds *DockerService) GetContainerLogs(containerName string, lines int) ([]string, error) {
	// Find container by name
	containers, err := ds.client.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var containerID string
	for _, container := range containers {
		name := strings.TrimPrefix(container.Names[0], "/")
		if name == containerName {
			containerID = container.ID
			break
		}
	}

	if containerID == "" {
		return nil, fmt.Errorf("container not found: %s", containerName)
	}

	// Get logs
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       strconv.Itoa(lines),
		Timestamps: true,
	}

	logs, err := ds.client.ContainerLogs(context.Background(), containerID, options)
	if err != nil {
		return nil, err
	}
	defer logs.Close()

	content, err := io.ReadAll(logs)
	if err != nil {
		return nil, err
	}

	// Parse logs (Docker logs have 8-byte header)
	var logLines []string
	lines_content := string(content)
	for _, line := range strings.Split(lines_content, "\n") {
		if len(line) > 8 {
			// Remove Docker log header
			logLines = append(logLines, line[8:])
		}
	}

	return logLines, nil
}

// RestartContainer restarts a container
func (ds *DockerService) RestartContainer(containerName string) error {
	containers, err := ds.client.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return err
	}

	var containerID string
	for _, container := range containers {
		name := strings.TrimPrefix(container.Names[0], "/")
		if name == containerName {
			containerID = container.ID
			break
		}
	}

	if containerID == "" {
		return fmt.Errorf("container not found: %s", containerName)
	}

	timeout := time.Second * 10
	return ds.client.ContainerRestart(context.Background(), containerID, &timeout)
}

// CheckUnhealthyContainers checks for unhealthy containers and performs auto-healing
func (ds *DockerService) CheckUnhealthyContainers() []models.AutoHealEvent {
	var events []models.AutoHealEvent

	// Check if auto-healing is enabled
	if !ds.config.AutoHeal.Enabled {
		return events
	}

	containers, err := ds.client.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		log.Printf("Error listing containers: %v", err)
		return events
	}

	for _, container := range containers {
		name := strings.TrimPrefix(container.Names[0], "/")
		
		// Check if container is in exclusion list
		excluded := false
		for _, excludedName := range ds.config.AutoHeal.ExcludeContainers {
			if name == excludedName {
				excluded = true
				break
			}
		}
		
		if excluded {
			continue
		}
		
		// Check if container is exited or unhealthy
		if container.State == "exited" || strings.Contains(container.Status, "unhealthy") {
			log.Printf("Found unhealthy container: %s (State: %s, Status: %s)", name, container.State, container.Status)
			
			// Attempt to restart
			err := ds.RestartContainer(name)
			success := err == nil
			
			event := models.AutoHealEvent{
				ContainerID: container.ID[:12],
				Name:        name,
				Action:      "restart",
				Reason:      fmt.Sprintf("Container state: %s", container.State),
				Success:     success,
				Timestamp:   time.Now(),
			}
			
			if !success {
				log.Printf("Failed to restart container %s: %v", name, err)
			} else {
				log.Printf("Successfully restarted container %s", name)
			}
			
			events = append(events, event)
		}
	}

	return events
}