package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/api/handlers"
)

// SetupSwaggerRoutes configures the Swagger documentation routes
func SetupSwaggerRoutes(e *echo.Echo) {
	handler := handlers.NewSwaggerHandler()

	// Swagger UI and spec routes
	e.GET("/swagger/*", handler.ServeUI)
	e.GET("/swagger.yaml", handler.ServeSpec)
	e.GET("/swagger", handler.RedirectToUI)
}
