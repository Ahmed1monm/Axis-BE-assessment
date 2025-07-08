package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/api/handlers"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/dtos"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/models"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of services.AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, input dtos.RegisterRequest) (*dtos.AuthResponse, error) {
	println("Register called with context ptr:", ctx)
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, input dtos.LoginRequest) (*dtos.AuthResponse, error) {
	println("Login called with context ptr:", ctx)
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.AuthResponse), args.Error(1)
}

func TestAuthHandler_Register(t *testing.T) {
	e := echo.New()
	mockAuthService := new(MockAuthService)
	handler := handlers.NewAuthHandler(mockAuthService)

	t.Run("Successful Registration", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := handlers.NewAuthHandler(mockAuthService)
		input := dtos.RegisterRequest{
			Name:        "John Doe",
			Email:       "john@example.com",
			Password:    "password123",
			PhoneNumber: "+1234567890",
		}

		response := &dtos.AuthResponse{
			Token: "jwt-token",
			User: &models.Account{
				Name:  input.Name,
				Email: input.Email,
			},
		}

		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockAuthService.On("Register", c.Request().Context(), input).Return(response, nil)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var responseBody dtos.AuthResponse
		err = json.Unmarshal(rec.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, response.Token, responseBody.Token)
		assert.Equal(t, response.User.Email, responseBody.User.Email)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		jsonBody := []byte(`{"invalid": json}`)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Validation Error", func(t *testing.T) {
		input := dtos.RegisterRequest{
			// Missing required fields
		}

		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Email Already Exists", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := handlers.NewAuthHandler(mockAuthService)
		input := dtos.RegisterRequest{
			Name:        "John Doe",
			Email:       "existing@example.com",
			Password:    "password123",
			PhoneNumber: "+1234567890",
		}

		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockAuthService.On("Register", c.Request().Context(), input).Return(nil, services.ErrEmailExists)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		input := dtos.RegisterRequest{
			Name:        "John Doe",
			Email:       "john@example.com",
			Password:    "password123",
			PhoneNumber: "+1234567890",
		}

		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockAuthService.On("Register", c.Request().Context(), input).Return((*dtos.AuthResponse)(nil), assert.AnError)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to register user", response["error"])
		mockAuthService.AssertExpectations(t)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	e := echo.New()
	mockAuthService := new(MockAuthService)
	handler := handlers.NewAuthHandler(mockAuthService)

	t.Run("Successful Login", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := handlers.NewAuthHandler(mockAuthService)
		input := dtos.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		response := &dtos.AuthResponse{
			Token: "jwt-token",
			User: &models.Account{
				Email: input.Email,
			},
		}

		mockAuthService.On("Login", c.Request().Context(), input).Return(response, nil)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var responseBody dtos.AuthResponse
		err = json.Unmarshal(rec.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, response.Token, responseBody.Token)
		assert.Equal(t, response.User.Email, responseBody.User.Email)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		jsonBody := []byte(`{"invalid": json}`)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Validation Error", func(t *testing.T) {
		input := dtos.LoginRequest{
			// Missing required fields
		}

		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := handlers.NewAuthHandler(mockAuthService)
		input := dtos.LoginRequest{
			Email:    "john@example.com",
			Password: "wrongpassword",
		}

		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockAuthService.On("Login", c.Request().Context(), input).Return(nil, services.ErrInvalidCredentials)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		input := dtos.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		jsonBody, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockAuthService.On("Login", c.Request().Context(), input).Return((*dtos.AuthResponse)(nil), assert.AnError)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to login", response["error"])
		mockAuthService.AssertExpectations(t)
	})
}
