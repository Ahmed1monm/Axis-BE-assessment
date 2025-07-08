package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/api/handlers"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/dtos"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/services"
	"github.com/Ahmed1monm/Axis-BE-assessment/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockTransactionServiceAdapter adapts a mock to be used as services.TransactionService
type MockTransactionServiceAdapter struct {
	Deposit     func(ctx context.Context, accountID primitive.ObjectID, amount float64, currency string) (string, error)
	Withdraw    func(ctx context.Context, accountID primitive.ObjectID, amount float64, currency string) (string, error)
	GetBalances func(ctx context.Context, accountID primitive.ObjectID) (*dtos.BalancesResponse, error)
}

func (m *MockTransactionServiceAdapter) Deposit(ctx context.Context, accountID primitive.ObjectID, amount float64, currency string) (string, error) {
	return m.Deposit(ctx, accountID, amount, currency)
}

func (m *MockTransactionServiceAdapter) Withdraw(ctx context.Context, accountID primitive.ObjectID, amount float64, currency string) (string, error) {
	return m.Withdraw(ctx, accountID, amount, currency)
}

func (m *MockTransactionServiceAdapter) GetBalances(ctx context.Context, accountID primitive.ObjectID) (*dtos.BalancesResponse, error) {
	return m.GetBalances(ctx, accountID)
}

func TestTransactionHandler_Deposit(t *testing.T) {
	e := echo.New()
	
	// Create mock service
	mockService := &services.TransactionService{}
	
	// Create adapter with mock implementations
	adapter := &MockTransactionServiceAdapter{}
	
	handler := handlers.NewTransactionHandler(mockService)

	t.Run("Successful Deposit", func(t *testing.T) {
		// Setup
		accountID := primitive.NewObjectID()
		input := dtos.TransactionRequest{
			AccountID: accountID.Hex(),
			Amount:    100.0,
			Currency:  "USD",
		}

		transactionID := primitive.NewObjectID().Hex()
		mockService.On("Deposit", mock.Anything, accountID, input.Amount, input.Currency).Return(transactionID, nil)

		// Create request
		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/transactions/deposit", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Test
		err := handler.Deposit(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response dtos.TransactionResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, transactionID, response.TransactionID)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		// Setup
		jsonBody := []byte(`{"invalid": json}`)
		req := httptest.NewRequest(http.MethodPost, "/transactions/deposit", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Test
		err := handler.Deposit(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Validation Error", func(t *testing.T) {
		// Setup - missing required fields
		input := dtos.TransactionRequest{
			// Missing AccountID, Amount, and Currency
		}

		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/transactions/deposit", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Test
		err := handler.Deposit(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "errors")
	})

	t.Run("Invalid Account ID", func(t *testing.T) {
		// Setup
		input := dtos.TransactionRequest{
			AccountID: "invalid-id",
			Amount:    100.0,
			Currency:  "USD",
		}

		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/transactions/deposit", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Test
		err := handler.Deposit(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid account ID", response["error"])
	})

	t.Run("Service Error - Custom Error", func(t *testing.T) {
		// Setup
		accountID := primitive.NewObjectID()
		input := dtos.TransactionRequest{
			AccountID: accountID.Hex(),
			Amount:    100.0,
			Currency:  "USD",
		}

		customErr := utils.NewError(http.StatusBadRequest, "insufficient funds")
		mockService.On("Deposit", mock.Anything, accountID, input.Amount, input.Currency).Return("", customErr)

		// Create request
		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/transactions/deposit", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Test
		err := handler.Deposit(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error - Generic Error", func(t *testing.T) {
		// Setup
		accountID := primitive.NewObjectID()
		input := dtos.TransactionRequest{
			AccountID: accountID.Hex(),
			Amount:    100.0,
			Currency:  "USD",
		}

		mockService.On("Deposit", mock.Anything, accountID, input.Amount, input.Currency).Return("", errors.New("database error"))

		// Create request
		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/transactions/deposit", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Test
		err := handler.Deposit(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestTransactionHandler_Withdraw(t *testing.T) {
	e := echo.New()
	
	// Create mock service
	mockService := &services.TransactionService{}
	
	// Create adapter with mock implementations
	adapter := &MockTransactionServiceAdapter{}
	
	handler := handlers.NewTransactionHandler(mockService)

	t.Run("Successful Withdrawal", func(t *testing.T) {
		// Setup
		accountID := primitive.NewObjectID()
		input := dtos.TransactionRequest{
			AccountID: accountID.Hex(),
			Amount:    50.0,
			Currency:  "USD",
		}

		transactionID := primitive.NewObjectID().Hex()
		mockService.On("Withdraw", mock.Anything, accountID, input.Amount, input.Currency).Return(transactionID, nil)

		// Create request
		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/transactions/withdraw", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Test
		err := handler.Withdraw(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response dtos.TransactionResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, transactionID, response.TransactionID)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		// Setup
		jsonBody := []byte(`{"invalid": json}`)
		req := httptest.NewRequest(http.MethodPost, "/transactions/withdraw", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Test
		err := handler.Withdraw(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Validation Error", func(t *testing.T) {
		// Setup - missing required fields
		input := dtos.TransactionRequest{
			// Missing AccountID, Amount, and Currency
		}

		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/transactions/withdraw", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Test
		err := handler.Withdraw(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "errors")
	})

	t.Run("Invalid Account ID", func(t *testing.T) {
		// Setup
		input := dtos.TransactionRequest{
			AccountID: "invalid-id",
			Amount:    50.0,
			Currency:  "USD",
		}

		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/transactions/withdraw", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Test
		err := handler.Withdraw(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid account ID", response["error"])
	})

	t.Run("Insufficient Funds Error", func(t *testing.T) {
		// Setup
		accountID := primitive.NewObjectID()
		input := dtos.TransactionRequest{
			AccountID: accountID.Hex(),
			Amount:    1000.0,
			Currency:  "USD",
		}

		customErr := utils.NewError(http.StatusBadRequest, "insufficient funds")
		mockService.On("Withdraw", mock.Anything, accountID, input.Amount, input.Currency).Return("", customErr)

		// Create request
		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/transactions/withdraw", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Test
		err := handler.Withdraw(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		// Setup
		accountID := primitive.NewObjectID()
		input := dtos.TransactionRequest{
			AccountID: accountID.Hex(),
			Amount:    50.0,
			Currency:  "USD",
		}

		mockService.On("Withdraw", mock.Anything, accountID, input.Amount, input.Currency).Return("", errors.New("database error"))

		// Create request
		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/transactions/withdraw", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Test
		err := handler.Withdraw(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestTransactionHandler_GetBalances(t *testing.T) {
	e := echo.New()
	
	// Create mock service
	mockService := &services.TransactionService{}
	
	// Create adapter with mock implementations
	adapter := &MockTransactionServiceAdapter{}
	
	handler := handlers.NewTransactionHandler(mockService)

	t.Run("Successful Get Balances", func(t *testing.T) {
		// Setup
		accountID := primitive.NewObjectID()
		
		balances := &dtos.BalancesResponse{
			Balances: []dtos.CurrencyBalance{
				{Currency: "USD", Amount: 100.0},
				{Currency: "EUR", Amount: 50.0},
			},
		}
		
		mockService.On("GetBalances", mock.Anything, accountID).Return(balances, nil)

		// Create request
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/transactions/balances/:account_id")
		c.SetParamNames("account_id")
		c.SetParamValues(accountID.Hex())

		// Test
		err := handler.GetBalances(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response dtos.BalancesResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.Balances))
		assert.Equal(t, "USD", response.Balances[0].Currency)
		assert.Equal(t, 100.0, response.Balances[0].Amount)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Account ID", func(t *testing.T) {
		// Setup
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/transactions/balances/:account_id")
		c.SetParamNames("account_id")
		c.SetParamValues("invalid-id")

		// Test
		err := handler.GetBalances(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Service Error - Custom Error", func(t *testing.T) {
		// Setup
		accountID := primitive.NewObjectID()
		customErr := utils.NewError(http.StatusNotFound, "account not found")
		
		mockService.On("GetBalances", mock.Anything, accountID).Return((*dtos.BalancesResponse)(nil), customErr)

		// Create request
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/transactions/balances/:account_id")
		c.SetParamNames("account_id")
		c.SetParamValues(accountID.Hex())

		// Test
		err := handler.GetBalances(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error - Generic Error", func(t *testing.T) {
		// Setup
		accountID := primitive.NewObjectID()
		
		mockService.On("GetBalances", mock.Anything, accountID).Return((*dtos.BalancesResponse)(nil), errors.New("database error"))

		// Create request
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/transactions/balances/:account_id")
		c.SetParamNames("account_id")
		c.SetParamValues(accountID.Hex())

		// Test
		err := handler.GetBalances(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}
