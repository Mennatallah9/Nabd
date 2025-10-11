package controllers

import (
	"net/http"
	"nabd/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AutoHealController struct {
	autoHealService *services.AutoHealService
}

// NewAutoHealController creates a new auto-heal controller
func NewAutoHealController(autoHealService *services.AutoHealService) *AutoHealController {
	return &AutoHealController{
		autoHealService: autoHealService,
	}
}

// GetAutoHealHistory returns auto-heal event history
func (ahc *AutoHealController) GetAutoHealHistory(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	events, err := ahc.autoHealService.GetAutoHealHistory(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": events})
}

// TriggerAutoHeal manually triggers auto-healing check
func (ahc *AutoHealController) TriggerAutoHeal(c *gin.Context) {
	ahc.autoHealService.PerformAutoHealing()
	c.JSON(http.StatusOK, gin.H{"message": "Auto-healing check triggered"})
}