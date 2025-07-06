package main

import (
	"context"

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

	// Initialize MongoDB with collections
	mongoClient, err := database.ConnectDB(cfg.MongoURI, cfg.DatabaseName)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Error().Err(err).Msg("Failed to disconnect from MongoDB")
		}
	}()

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
