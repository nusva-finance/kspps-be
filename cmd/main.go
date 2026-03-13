package main

import (
	"log"
	"nusvakspps/config"
	"nusvakspps/middleware"
	"nusvakspps/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables from .env file
	if err := middleware.LoadEnv(); err != nil {
		log.Println("Warning: .env file not found, using default environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create router
	router := gin.Default()

	// Apply middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.RequestLoggerMiddleware())

	// Register routes
	routes.RegisterRoutes(router)

	// Start server
	addr := ":" + cfg.Server.Port
	log.Printf("Server starting on port %s...", cfg.Server.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}