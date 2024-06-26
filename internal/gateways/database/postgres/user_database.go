package postgres

import (
	"context"
	"strings"
	"weather-notification/internal/domain"
	"weather-notification/internal/domain/entities"
	"weather-notification/internal/domain/ports"
	"weather-notification/internal/gateways/database"
	"weather-notification/internal/gateways/database/postgres/models"
)

var _ ports.UserDatabaseGateway = (*userDatabase)(nil)

type userDatabase struct {
	Client *Client
}

func NewUserDatabase(client *Client) *userDatabase {
	return &userDatabase{
		Client: client,
	}
}

func (u *userDatabase) InsertUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	model := models.NewUserModel(user)

	_, err := u.Client.DB.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), database.DuplicateKeyPrefix) {
			return &entities.User{}, domain.ErrEmailIsAlreadyInUse
		}

		return &entities.User{}, database.NewInsertError(models.UsersTableName, err)
	}

	return model.ToEntity(), nil
}

func (g *userDatabase) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	return nil, nil
}
