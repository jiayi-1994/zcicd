package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
	Validate.RegisterValidation("not_blank", notBlank)
	Validate.RegisterValidation("safe_string", safeString)
}

func notBlank(fl validator.FieldLevel) bool {
	return strings.TrimSpace(fl.Field().String()) != ""
}

func safeString(fl validator.FieldLevel) bool {
	dangerous := []string{"'", "\"", ";", "--", "/*", "*/", "xp_", "DROP ", "SELECT ", "INSERT ", "DELETE ", "UPDATE "}
	val := strings.ToUpper(fl.Field().String())
	for _, d := range dangerous {
		if strings.Contains(val, strings.ToUpper(d)) {
			return false
		}
	}
	return true
}

func ValidateStruct(s interface{}) error {
	return Validate.Struct(s)
}

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errors[e.Field()] = fmt.Sprintf("%s is required", e.Field())
			case "not_blank":
				errors[e.Field()] = fmt.Sprintf("%s must not be blank", e.Field())
			case "safe_string":
				errors[e.Field()] = fmt.Sprintf("%s contains invalid characters", e.Field())
			default:
				errors[e.Field()] = fmt.Sprintf("%s failed on %s validation", e.Field(), e.Tag())
			}
		}
	}
	return errors
}
