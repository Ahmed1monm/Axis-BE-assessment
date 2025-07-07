package routes

import (
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/api/handlers"
	"github.com/labstack/echo/v4"
)

// SetupTransactionRoutes sets up all transaction related routes
// @Summary Setup transaction routes
// @Description Configures deposit and withdrawal endpoints under /api/v1/transactions
// @Tags transactions
func SetupTransactionRoutes(g *echo.Group, h *handlers.TransactionHandler) {
	transactions := g.Group("/transactions")

	// POST /api/v1/transactions/deposit
	transactions.POST("/deposit", h.Deposit)

	// POST /api/v1/transactions/withdraw
	transactions.POST("/withdraw", h.Withdraw)
}
