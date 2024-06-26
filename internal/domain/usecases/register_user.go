package usecases

import (
	"context"
	"weather-notification/internal/domain/entities"
	"weather-notification/internal/domain/ports"
)

var _ RegisterUserUseCase = (*registerUserUseCase)(nil)

type RegisterUserUseCase interface {
	Execute(ctx context.Context, user *entities.User) (*entities.User, error)
}

type registerUserUseCase struct {
	userDatabase ports.UserDatabaseGateway
}

func NewRegisterUseCase(userDatabase ports.UserDatabaseGateway) *registerUserUseCase {
	return &registerUserUseCase{
		userDatabase: userDatabase,
	}
}

func (u *registerUserUseCase) Execute(ctx context.Context, user *entities.User) (*entities.User, error) {
	err := user.Validate()
	if err != nil {
		return nil, err
	}

	createdUser, err := u.userDatabase.InsertUser(ctx, user)
	if err != nil {
		return &entities.User{}, err
	}

	return createdUser, nil
}
