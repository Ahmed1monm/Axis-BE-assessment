package utils

import "os"

// GetEnv retrieves an environment variable value or returns a default value if not set
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
