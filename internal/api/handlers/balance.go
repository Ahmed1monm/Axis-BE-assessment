package handlers

import (
	"net/http"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/services"
	"github.com/Ahmed1monm/Axis-BE-assessment/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BalanceHandler struct {
	balanceService *services.BalanceService
}

func NewBalanceHandler(balanceService *services.BalanceService) *BalanceHandler {
	return &BalanceHandler{
		balanceService: balanceService,
	}
}

// GetBalances handles the GET /balances/:account_id endpoint
func (h *BalanceHandler) GetBalances(c echo.Context) error {
	accountIDStr := c.Param("account_id")
	accountID, err := primitive.ObjectIDFromHex(accountIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(
			http.StatusBadRequest,
			"invalid account ID",
		))
	}

	response, err := h.balanceService.GetBalances(c.Request().Context(), accountID)
	if err != nil {
		if customErr, ok := utils.IsCustomError(err); ok {
			return c.JSON(customErr.Code, customErr)
		}
		return c.JSON(http.StatusInternalServerError, utils.WrapError(err, http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, response)
}
