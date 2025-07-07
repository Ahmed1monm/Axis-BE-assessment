package dtos

import "github.com/Ahmed1monm/Axis-BE-assessment/internal/models"

// RegisterRequest represents the registration request data
type RegisterRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
	Password    string `json:"password" validate:"required,min=8,max=72"`
}

// LoginRequest represents the login request data
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string          `json:"token"`
	User  *models.Account `json:"user"`
}
