package routes

import (
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/api/handlers"
	"github.com/labstack/echo/v4"
)

// SetupBalanceRoutes sets up all balance related routes
// @Summary Setup balance routes
// @Description Configures balance endpoints under /api/v1/balances
// @Tags balances
func SetupBalanceRoutes(g *echo.Group, h *handlers.BalanceHandler) {
	balances := g.Group("/balances")

	// GET /api/v1/balances/:account_id
	balances.GET("/:account_id", h.GetBalances)
}
