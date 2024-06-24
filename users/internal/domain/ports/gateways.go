package ports

import (
	"context"
	"weather-notification/users/internal/domain/entities"
)

type UserDatabaseGateway interface {
	InsertUser(ctx context.Context, user *entities.User) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error)
}
