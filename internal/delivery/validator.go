package delivery

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type PlaygroundValidator struct {
	validator *validator.Validate
}

func (v *PlaygroundValidator) formatErrorMsg(err *validator.FieldError) string {
	return fmt.Sprintf(
		"Field validation for \"%s\" failed on the \"%s\" tag",
		(*err).Field(),
		(*err).Tag(),
	)
}

func (v *PlaygroundValidator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	if err == nil {
		return nil
	}

	errorMessages := []string{}
	for _, err := range err.(validator.ValidationErrors) {
		errorMessages = append(
			errorMessages,
			v.formatErrorMsg(&err),
		)
	}
	message := strings.Join(errorMessages, "\n")

	return errors.New(message)
}

func NewValidator() *PlaygroundValidator {
	return &PlaygroundValidator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
