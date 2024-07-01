package entities

import (
	"time"
)

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Phone     string
	CityID    int
	OptIn     bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(firstName, lastName, email, phone string, cityID int) *User {
	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		CityID:    cityID,
		OptIn:     true,
	}
}

func (u *User) Subscribe() {
	u.OptIn = true
}

func (u *User) Unsubscribe() {
	u.OptIn = false
}
