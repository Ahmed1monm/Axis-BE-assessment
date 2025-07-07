package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/api/handlers"
)

func SetupAuthRoutes(e *echo.Echo, authHandler *handlers.AuthHandler) {
	auth := e.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
}
