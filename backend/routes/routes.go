package routes

import (
	"nabd/controllers"
	"nabd/models"
	"nabd/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(
	containerController *controllers.ContainerController,
	autoHealController *controllers.AutoHealController,
	alertController *controllers.AlertController,
	authController *controllers.AuthController,
	config *models.Config,
) *gin.Engine {
	
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)
	
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Serve static files (frontend)
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	
	// Serve index.html for all non-API routes (SPA support)
	router.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.File("./static/index.html")
		} else {
			c.JSON(404, gin.H{"error": "Not found"})
		}
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "version": "v0.1.0"})
	})

	// Authentication routes (no auth required)
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", authController.Login)
	}

	// API routes (auth required)
	api := router.Group("/api")
	api.Use(utils.AuthMiddleware(config.Auth.AdminToken))
	{
		// Container routes
		api.GET("/containers", containerController.GetContainers)
		api.GET("/metrics", containerController.GetMetrics)
		api.GET("/metrics/:id/history", containerController.GetMetricsHistory)
		api.GET("/logs", containerController.GetLogs)
		api.POST("/containers/:name/restart", containerController.RestartContainer)

		// Auto-heal routes
		api.GET("/autoheal/history", autoHealController.GetAutoHealHistory)
		api.POST("/autoheal/trigger", autoHealController.TriggerAutoHeal)

		// Alert routes
		api.GET("/alerts", alertController.GetAlerts)
	}

	return router
}