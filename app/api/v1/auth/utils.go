package auth

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func getAuthErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return fmt.Sprintf("Invalid email %s provided", fieldError.Value())
	}
	return "Unknown error"
}
