package handlers

import (
	"net/http"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/api/validation"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/dtos"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/services"
	"github.com/Ahmed1monm/Axis-BE-assessment/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

// GetBalances handles the GET /transactions/balances/:account_id endpoint
func (h *TransactionHandler) GetBalances(c echo.Context) error {
	accountIDStr := c.Param("account_id")
	accountID, err := primitive.ObjectIDFromHex(accountIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(
			http.StatusBadRequest,
			"invalid account ID",
		))
	}

	response, err := h.transactionService.GetBalances(c.Request().Context(), accountID)
	if err != nil {
		if customErr, ok := utils.IsCustomError(err); ok {
			return c.JSON(customErr.Code, customErr)
		}
		return c.JSON(http.StatusInternalServerError, utils.WrapError(err, http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, response)
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (h *TransactionHandler) Deposit(c echo.Context) error {
	var input dtos.TransactionRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if errors := validation.ValidateStruct(input); len(errors) > 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"errors": errors})
	}

	accountID, err := primitive.ObjectIDFromHex(input.AccountID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid account ID"})
	}

	transactionID, err := h.transactionService.Deposit(c.Request().Context(), accountID, input.Amount, input.Currency)
	if err != nil {
		if customErr, ok := utils.IsCustomError(err); ok {
			return c.JSON(customErr.Code, customErr)
		}
		return c.JSON(http.StatusInternalServerError, utils.WrapError(err, http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, &dtos.TransactionResponse{
		TransactionID: transactionID,
	})
}

func (h *TransactionHandler) Withdraw(c echo.Context) error {
	var input dtos.TransactionRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if errors := validation.ValidateStruct(input); len(errors) > 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"errors": errors})
	}

	accountID, err := primitive.ObjectIDFromHex(input.AccountID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid account ID"})
	}

	transactionID, err := h.transactionService.Withdraw(c.Request().Context(), accountID, input.Amount, input.Currency)
	if err != nil {
		if customErr, ok := utils.IsCustomError(err); ok {
			return c.JSON(customErr.Code, customErr)
		}
		return c.JSON(http.StatusInternalServerError, utils.WrapError(err, http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, &dtos.TransactionResponse{
		TransactionID: transactionID,
	})
}
