package entities

import "github.com/go-playground/validator/v10"

type User struct {
	ID        int64
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Email     string `validate:"required,email"`
	Phone     string `validate:"required"`
	OptIn     bool
}

func NewUser(firstName, lastName, email, phone string, optIn bool) *User {
	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		OptIn:     optIn,
	}
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
