package utilities

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(s interface{}) string {
	var errors []string
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf(
				"%s is %s",
				err.Field(), err.Tag(),
			))
		}
		return strings.Join(errors, ", ")
	}
	return ""
}
