package routes

import (
	"axis-be-assessment/internal/api/handlers"
	"axis-be-assessment/internal/api/middleware"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(e *echo.Echo, db *mongo.Client, logger zerolog.Logger) {
	// Middleware
	e.Use(echomw.Recover())
	e.Use(echomw.CORS())
	e.Use(middleware.RequestLogger(logger))

	// Health Check
	e.GET("/health", handlers.HealthCheck())

	// API v1 group
	v1 := e.Group("/api/v1")
	
	// Add your routes here
	// Example: v1.POST("/users", handlers.CreateUser(db))
}
