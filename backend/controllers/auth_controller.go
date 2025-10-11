package controllers

import (
	"net/http"
	"nabd/models"
	"nabd/services"
	"nabd/utils"

	"github.com/gin-gonic/gin"
)

type AlertController struct {
	metricsService *services.MetricsService
}

// NewAlertController creates a new alert controller
func NewAlertController(metricsService *services.MetricsService) *AlertController {
	return &AlertController{
		metricsService: metricsService,
	}
}

// GetAlerts returns active alerts
func (ac *AlertController) GetAlerts(c *gin.Context) {
	alerts, err := ac.metricsService.GetActiveAlerts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": alerts})
}

// AuthController handles authentication
type AuthController struct {
	config *models.Config
}

// NewAuthController creates a new auth controller
func NewAuthController(config *models.Config) *AuthController {
	return &AuthController{
		config: config,
	}
}

// Login handles user authentication
func (ac *AuthController) Login(c *gin.Context) {
	var loginRequest struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if loginRequest.Token != ac.config.Auth.AdminToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid admin token"})
		return
	}

	// Generate JWT token
	jwtToken, err := utils.GenerateToken(ac.config.Auth.AdminToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": jwtToken,
		"message": "Authentication successful",
	})
}