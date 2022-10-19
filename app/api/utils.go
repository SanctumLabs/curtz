package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/sanctumlabs/curtz/app/pkg"
)

func ParseDtoFieldError(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return fmt.Sprintf("Invalid email %s provided", fieldError.Value())
	case "expires_on":
		return fmt.Sprintf("Invalid Expires On field %s provided, expected format %s", fieldError.Value(), pkg.DateFormat)
	default:
		return "Unknown error"
	}
}
