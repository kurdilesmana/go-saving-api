package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type RequestValidator struct {
	validator *validator.Validate
}

func NewRequestValidator() *RequestValidator {
	return &RequestValidator{validator: validator.New()}
}

func (v *RequestValidator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		var message string
		var fieldRequired []string
		var fieldOneof []string

		if castedObject, ok := err.(validator.ValidationErrors); ok {
			for _, err := range castedObject {
				switch err.Tag() {
				case "required":
					field := fmt.Sprintf("%v,", toSnakeCase(err.Field()))
					fieldRequired = append(fieldRequired, field)
				case "oneof":
					field := toSnakeCase(err.Field())
					value := err.Param()
					oneOf := fmt.Sprintf("%v: %v,", field, value)
					fieldOneof = append(fieldOneof, oneOf)
				}
			}

			if len(fieldRequired) > 0 {
				message = fmt.Sprintf("REQUIRED: field %v harus diisi. ", strings.Join(fieldRequired, ","))
			}

			if len(fieldOneof) > 0 {
				message = message + fmt.Sprintf("ENUM: %v. ", fieldOneof)
			}
		}
		return fmt.Errorf(message)
	}
	return nil
}

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
