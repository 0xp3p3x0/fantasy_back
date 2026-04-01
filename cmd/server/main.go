package main

import (
	_ "back/docs"
	"back/internal/config"
	"back/internal/db"
	"back/internal/logger"
	"back/internal/server"
	"log"

	"github.com/joho/godotenv"
)

// @title Fantasy Back API
// @version 1.0
// @description Backend API for authentication, profile, casino and agent management.
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
		// Don't fatal here, as .env might not exist in production
	}

	// Initialize logger
	if err := logger.Init(); err != nil {
		logger.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	database := db.Init(cfg)

	// Create and start server
	srv, err := server.NewServer(database, cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
		logger.Fatalf("Failed to create server: %v", err)
	}

	if err := srv.Run(cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		logger.Fatalf("Failed to start server: %v", err)
	}
}
