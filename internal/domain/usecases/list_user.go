package usecases

import (
	"context"
	"fmt"
	"weather-notification/internal/domain/entities"
	"weather-notification/internal/domain/ports"
)

var _ ListUserUseCase = (*listUserUseCase)(nil)

type ListUserUseCase interface {
	Execute(ctx context.Context) ([]*entities.User, error)
}

type listUserUseCase struct {
	userDatabase ports.UserDatabaseGateway
}

func NewListUserUseUseCase(userDatabase ports.UserDatabaseGateway) *listUserUseCase {
	return &listUserUseCase{
		userDatabase: userDatabase,
	}
}

func (u *listUserUseCase) Execute(ctx context.Context) ([]*entities.User, error) {
	users, err := u.userDatabase.ListUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("unsubscribe usecase - failed to get user from database: %w", err)
	}

	return users, nil
}
