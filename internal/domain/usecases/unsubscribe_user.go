package usecases

import (
	"context"
	"fmt"
	"weather-notification/internal/domain/entities"
	"weather-notification/internal/domain/ports"
)

var _ UnsubscribeUserUseCase = (*unsubscribeUserUseCase)(nil)

type UnsubscribeUserUseCase interface {
	Execute(ctx context.Context, email string) (*entities.User, error)
}

type unsubscribeUserUseCase struct {
	userDatabase ports.UserDatabaseGateway
}

func NewUnsubscribeUseUseCase(userDatabase ports.UserDatabaseGateway) *unsubscribeUserUseCase {
	return &unsubscribeUserUseCase{
		userDatabase: userDatabase,
	}
}

func (u *unsubscribeUserUseCase) Execute(ctx context.Context, email string) (*entities.User, error) {
	user, err := u.userDatabase.GetUser(ctx, email)
	if err != nil {
		return &entities.User{}, fmt.Errorf("unsubscribe usecase - failed to get user from database: %w", err)
	}

	user.Unsubscribe()
	updatedUser, err := u.userDatabase.UpdateUser(ctx, user)
	if err != nil {
		return &entities.User{}, fmt.Errorf("unsubscribe usecase - failed to update user on database: %w", err)
	}

	return updatedUser, nil
}
