package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/dtos"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	// Add custom validations here if needed
}

// ValidateStruct validates a struct using validator tags and returns structured validation errors
func ValidateStruct(s interface{}) dtos.ValidationErrors {
	var errors dtos.ValidationErrors

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, dtos.ValidationError{
				Field:   strings.ToLower(err.Field()),
				Message: generateValidationMessage(err),
			})
		}
	}

	return errors
}

// generateValidationMessage generates a human-readable validation message
func generateValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", err.Field(), err.Param())
	case "e164":
		return "Invalid phone number format. Must be in E.164 format"
	default:
		return fmt.Sprintf("%s is not valid", err.Field())
	}
}
