package main

import (
	"context"
	"log"

	"axis-be-assessment/internal/api/routes"
	"axis-be-assessment/internal/config"
	"axis-be-assessment/pkg/database"
	"axis-be-assessment/pkg/logger"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize logger
	log := logger.New()

	// Load configuration
	cfg := config.Load()

	// Initialize MongoDB
	mongoClient, err := database.ConnectDB(cfg.MongoURI)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}
	defer mongoClient.Disconnect(context.Background())

	// Initialize Echo
	e := echo.New()

	// Setup routes
	routes.Setup(e, mongoClient, log)

	// Start server
	log.Info().Msgf("Server starting on port %s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
