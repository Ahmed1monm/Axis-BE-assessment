package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealthCheck() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "OK",
		})
	}
}
