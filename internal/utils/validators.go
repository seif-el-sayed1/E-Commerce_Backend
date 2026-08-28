package utils

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FormatValidationError extracts a user-friendly error message from a binding error
func FormatValidationError(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		fe := ve[0]
		field := CamelToWords(fe.Field())
		
		switch fe.Tag() {
		case "required":
			return field + " is required"
		case "email":
			return "Invalid email format"
		case "min":
			return field + " must have at least " + fe.Param() + " characters"
		case "max":
			return field + " must have at most " + fe.Param() + " characters"
		case "gte":
			return field + " must be greater than or equal to " + fe.Param()
		case "lte":
			return field + " must be less than or equal to " + fe.Param()
		case "gt":
			return field + " must be greater than " + fe.Param()
		case "lt":
			return field + " must be less than " + fe.Param()
		case "uuid":
			return field + " must be a valid UUID"
		default:
			return "Invalid value for " + field
		}
	}

	// For JSON unmarshaling errors, etc (e.g., passing a string instead of an int)
	if strings.Contains(err.Error(), "Unmarshal type error") || strings.Contains(err.Error(), "cannot unmarshal") {
		return "Invalid data type in request body"
	}

	return "Invalid request body"
}
