package ports

import (
	"context"
	"weather-notification/internal/domain/entities"
)

type UserDatabaseGateway interface {
	GetUser(ctx context.Context, email string) (*entities.User, error)
	InsertUser(ctx context.Context, user *entities.User) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error)
}
