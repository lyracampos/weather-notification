package usecases

import (
	"context"
	"weather-notification/users/internal/domain/entities"
	"weather-notification/users/internal/domain/ports"
)

type RegisterUseCase struct {
	databaseGateway ports.UserDatabaseGateway
}

func NewRegisterUseCase(databaseGateway ports.UserDatabaseGateway) *RegisterUseCase {
	return &RegisterUseCase{
		databaseGateway: databaseGateway,
	}
}

func (u *RegisterUseCase) Execute(ctx context.Context, user *entities.User) (*entities.User, error) {
	err := user.Validate()
	if err != nil {
		return nil, err
	}

	createdUser, err := u.databaseGateway.InsertUser(ctx, user)
	if err != nil {
		return &entities.User{}, err
	}

	return createdUser, nil
}
