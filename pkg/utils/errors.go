package utils

import (
	"fmt"
	"net/http"
)

// CustomError represents a custom error with HTTP status code and message
type CustomError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
}

// Error implements the error interface
func (e *CustomError) Error() string {
	return e.Message
}

// NewError creates a new CustomError
func NewError(code int, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

// Common application errors
var (
	ErrInvalidAmount = NewError(
		http.StatusBadRequest,
		"invalid amount: must be greater than 0",
	)

	ErrInsufficientBalance = NewError(
		http.StatusBadRequest,
		"insufficient balance",
	)

	ErrInvalidCredentials = NewError(
		http.StatusUnauthorized,
		"invalid credentials",
	)

	ErrUserNotFound = NewError(
		http.StatusNotFound,
		"user not found",
	)

	ErrUserAlreadyExists = NewError(
		http.StatusConflict,
		"user already exists",
	)

	ErrInvalidToken = NewError(
		http.StatusUnauthorized,
		"invalid or expired token",
	)
)

// IsCustomError checks if an error is a CustomError
func IsCustomError(err error) (*CustomError, bool) {
	customErr, ok := err.(*CustomError)
	return customErr, ok
}

// WrapError wraps a standard error into a CustomError with a specific HTTP status code
func WrapError(err error, code int) *CustomError {
	return NewError(code, err.Error())
}

// DatabaseError wraps database-related errors
func DatabaseError(operation string, err error) *CustomError {
	return NewError(
		http.StatusInternalServerError,
		fmt.Sprintf("database error during %s: %v", operation, err),
	)
}
