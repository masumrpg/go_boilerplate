package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the go-playground/validator
type Validator struct {
	validate *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// ValidateStruct validates a struct
func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}

// GetValidationErrors returns formatted validation errors
func GetValidationErrors(err error) []string {
	var errors []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, formatValidationError(e))
		}
	} else {
		errors = append(errors, err.Error())
	}

	return errors
}

// formatValidationError formats a single validation error
func formatValidationError(e validator.FieldError) string {
	field := e.Field()
	tag := e.Tag()
	param := e.Param()

	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email"
	case "min":
		return field + " must be at least " + param + " characters"
	case "max":
		return field + " must be at most " + param + " characters"
	case "len":
		return field + " must be " + param + " characters"
	default:
		return field + " failed on " + tag + " validation"
	}
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
