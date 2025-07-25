package routes

import (
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/api/handlers"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/api/middleware"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/services"
)

func Setup(e *echo.Echo, db *mongo.Client, logger zerolog.Logger) {
	// Middleware
	e.Use(echomw.Recover())
	e.Use(echomw.CORS())
	e.Use(middleware.RequestLogger(logger))

	// Health Check
	e.GET("/health", handlers.HealthCheck())

	// Setup Swagger documentation routes
	SetupSwaggerRoutes(e)

	// API v1 group
	v1 := e.Group("/api/v1")

	// Public routes (no authentication required)
	SetupAuthRoutes(v1, db)

	// Protected routes (authentication required)
	protected := v1.Group("", middleware.Auth())

	// Transaction routes
	transactionHandler := handlers.NewTransactionHandler(services.NewTransactionService(db))
	SetupTransactionRoutes(protected, transactionHandler)

	// Balance routes
	balanceHandler := handlers.NewBalanceHandler(services.NewBalanceService(db))
	SetupBalanceRoutes(protected, balanceHandler)
}
