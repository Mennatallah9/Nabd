package controllers

import (
	"net/http"
	"nabd/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContainerController struct {
	dockerService  *services.DockerService
	metricsService *services.MetricsService
}

// NewContainerController creates a new container controller
func NewContainerController(dockerService *services.DockerService, metricsService *services.MetricsService) *ContainerController {
	return &ContainerController{
		dockerService:  dockerService,
		metricsService: metricsService,
	}
}

// GetContainers returns all containers
func (cc *ContainerController) GetContainers(c *gin.Context) {
	containers, err := cc.dockerService.GetContainers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": containers})
}

// GetMetrics returns current metrics for all containers
func (cc *ContainerController) GetMetrics(c *gin.Context) {
	metrics, err := cc.metricsService.GetLatestMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": metrics})
}

// GetMetricsHistory returns historical metrics for a container
func (cc *ContainerController) GetMetricsHistory(c *gin.Context) {
	containerID := c.Param("id")
	hoursStr := c.DefaultQuery("hours", "24")
	
	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hours parameter"})
		return
	}

	metrics, err := cc.metricsService.GetMetricsHistory(containerID, hours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": metrics})
}

// GetLogs returns logs for a specific container
func (cc *ContainerController) GetLogs(c *gin.Context) {
	containerName := c.Query("container")
	if containerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "container parameter is required"})
		return
	}

	linesStr := c.DefaultQuery("lines", "100")
	lines, err := strconv.Atoi(linesStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lines parameter"})
		return
	}

	logs, err := cc.dockerService.GetContainerLogs(containerName, lines)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

// RestartContainer restarts a specific container
func (cc *ContainerController) RestartContainer(c *gin.Context) {
	containerName := c.Param("name")
	
	err := cc.dockerService.RestartContainer(containerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container restarted successfully"})
}