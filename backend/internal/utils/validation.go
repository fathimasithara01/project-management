package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) string {
	var errors []string

	// type cast
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error() // fallback
	}

	for _, e := range validationErrors {
		field := strings.ToLower(e.Field())

		switch e.Tag() {

		case "required":
			errors = append(errors, field+" is required")

		case "email":
			errors = append(errors, field+" must be a valid email")

		case "min":
			errors = append(errors, field+" must be at least "+e.Param()+" characters")

		case "max":
			errors = append(errors, field+" must be at most "+e.Param()+" characters")

		case "oneof":
			errors = append(errors, field+" must be one of ["+e.Param()+"]")

		default:
			errors = append(errors, field+" is invalid")
		}
	}

	return strings.Join(errors, ", ")
}