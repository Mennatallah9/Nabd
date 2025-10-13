package main

import (
	"log"
	"nabd/controllers"
	"nabd/routes"
	"nabd/services"
	"nabd/utils"
	"time"
)

func main() {
	log.Println("Starting Nabd - Container Observability & Auto-Healing Tool")

	// Load configuration
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	if err := utils.InitDatabase(config.Database.Path); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Docker service
	dockerService, err := services.NewDockerService(config)
	if err != nil {
		log.Fatalf("Failed to initialize Docker service: %v", err)
	}

	// Initialize metrics service
	metricsService := services.NewMetricsService(dockerService, config)

	// Initialize auto-heal service
	autoHealService := services.NewAutoHealService(dockerService, metricsService, config)

	// Start background services
	autoHealService.StartAutoHealing()

	// Start metrics collection
	go func() {
		ticker := time.NewTicker(15 * time.Second) // Collect metrics every 15 seconds
		for range ticker.C {
			if err := metricsService.CollectAndStoreMetrics(); err != nil {
				log.Printf("Error collecting metrics: %v", err)
			}
		}
	}()

	// Initialize controllers
	containerController := controllers.NewContainerController(dockerService, metricsService)
	autoHealController := controllers.NewAutoHealController(autoHealService)
	alertController := controllers.NewAlertController(metricsService)
	authController := controllers.NewAuthController(config)

	// Setup routes
	router := routes.SetupRoutes(
		containerController,
		autoHealController,
		alertController,
		authController,
		config,
	)

	log.Println("Nabd server started on port 8080")
	log.Printf("Admin token: %s", config.Auth.AdminToken)
	
	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}