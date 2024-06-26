package entities

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// swagger:model
type User struct {
	// user identification ID
	//
	ID int64
	// the user's first name
	//
	// required: true
	FirstName string `validate:"required" json:"first_name"`
	// the user's last name
	//
	// required: true
	LastName string `validate:"required" json:"last_name"`
	// the user's email (unique)
	//
	// required: true
	Email string `validate:"required,email"`
	// the user's phone
	//
	// required: true
	Phone string `validate:"required"`
	// the option to receive user notifications
	//
	// required: true
	OptIn bool

	CreatedAt time.Time
	UpdatedAt time.Time
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

func (u *User) Subscribe() {
	u.OptIn = true
}

func (u *User) Unsubscribe() {
	u.OptIn = false
}
