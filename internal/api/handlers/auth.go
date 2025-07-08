package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/api/validation"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/dtos"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/services"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c echo.Context) error {
	var input dtos.RegisterRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if errors := validation.ValidateStruct(input); len(errors) > 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"errors": errors})
	}

	response, err := h.authService.Register(c.Request().Context(), input)
	fmt.Printf("DEBUG: Register response=%%#v, err=%%v\n", response, err)
	if err != nil {
		if err == services.ErrEmailExists {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Email already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to register user"})
	}

	return c.JSON(http.StatusCreated, response)
}

// Login handles user login
func (h *AuthHandler) Login(c echo.Context) error {
	var input dtos.LoginRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if errors := validation.ValidateStruct(input); len(errors) > 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"errors": errors})
	}

	response, err := h.authService.Login(c.Request().Context(), input)
	fmt.Printf("DEBUG: Login response=%%#v, err=%%v\n", response, err)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to login"})
	}

	return c.JSON(http.StatusOK, response)
}
