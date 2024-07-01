package usecases

import (
	"context"
	"fmt"
	"weather-notification/internal/domain/entities"
	"weather-notification/internal/domain/ports"

	"github.com/go-playground/validator"
)

var _ RegisterUserUseCase = (*registerUserUseCase)(nil)

// swagger:model
type RegisterUserInput struct {
	// the user's first name
	//
	// required: true
	FirstName string `validate:"required" json:"first_name"`
	// the user's last name
	//
	// required: true
	LastName string `validate:"required" json:"last_name"`
	// the email which user will be notified
	//
	// required: true
	Email string `validate:"required,email"`
	// the phone which user will be notified
	//
	// required: true
	Phone string `validate:"required"`
	// the city where user is from
	//
	// required: true
	City string `validate:"required"`
	// the state where user is from
	//
	// required: true
	State string `validate:"required"`
}

func (u *RegisterUserInput) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

type RegisterUserUseCase interface {
	Execute(ctx context.Context, input *RegisterUserInput) (*entities.User, error)
}

type registerUserUseCase struct {
	userDatabase ports.UserDatabaseGateway
	weatherAPI   ports.WeatherHTTPGateway
}

func NewRegisterUseCase(userDatabase ports.UserDatabaseGateway, weatherAPI ports.WeatherHTTPGateway) *registerUserUseCase {
	return &registerUserUseCase{
		userDatabase: userDatabase,
		weatherAPI:   weatherAPI,
	}
}

func (u *registerUserUseCase) Execute(ctx context.Context, input *RegisterUserInput) (*entities.User, error) {
	err := input.Validate()
	if err != nil {
		return &entities.User{}, fmt.Errorf("registerUser usecase - invalid request: %w", err)
	}

	//TODO: get city from cache

	city, err := u.weatherAPI.GetCity(ctx, input.City, input.State)
	if err != nil {
		return &entities.User{}, fmt.Errorf("registerUser usecase - failed to get city from weather API: %w", err)
	}

	user := entities.NewUser(input.FirstName, input.LastName, input.Email, input.Phone, city.ID)
	createdUser, err := u.userDatabase.InsertUser(ctx, user)
	if err != nil {
		return &entities.User{}, fmt.Errorf("registerUser usecase - failed to insert an user into database: %w", err)
	}

	return createdUser, nil
}
