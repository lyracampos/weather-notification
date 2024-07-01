package domain

import (
	"errors"
)

var (
	ErrEmailIsAlreadyInUse = errors.New("email is arealdy registered for another user")
	ErrUserNotFound        = errors.New("user not found")
	ErrCityNotFound        = errors.New("city not found")
)
