package domain

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var (
	ErrEmailIsAlreadyInUse = errors.New("email is arealdy in use")
	ErrUserNotFound        = errors.New("user not found")
)

func Error(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf("the field '%s' should not be empty", fieldError.Field())
	default:
		return fmt.Sprintf("the field '%s' is invalid", fieldError.Field())
	}
}
