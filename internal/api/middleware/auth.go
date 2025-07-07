package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"axis-be-assessment/pkg/jwt"
)

// UserIDKey is the key used to store the user ID in the context
const UserIDKey = "user_id"

// Auth returns a middleware function that authenticates requests using JWT
func Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header is required",
				})
			}

			// Check if the Authorization header has the correct format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header must be in the format: Bearer {token}",
				})
			}

			// Extract and validate the token
			tokenString := parts[1]
			claims, err := jwt.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid or expired token",
				})
			}

			// Add user ID to context
			c.Set(UserIDKey, claims.UserID)

			return next(c)
		}
	}
}

// GetUserID retrieves the authenticated user's ID from the context
func GetUserID(c echo.Context) uint {
	userID, ok := c.Get(UserIDKey).(uint)
	if !ok {
		return 0 // Return 0 if user ID is not found or invalid
	}
	return userID
}
