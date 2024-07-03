package postgres

import (
	"context"
	"fmt"
	"strings"
	"weather-notification/internal/domain"
	"weather-notification/internal/domain/entities"
	"weather-notification/internal/domain/ports"
	"weather-notification/internal/gateways/database"
	"weather-notification/internal/gateways/database/postgres/models"

	"github.com/uptrace/bun"
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

func (u *userDatabase) GetUser(ctx context.Context, email string) (*entities.User, error) {
	model := models.User{}

	if err := u.Client.DB.NewSelect().Model(&model).Where("? = ?", bun.Ident("email"), email).Scan(ctx); err != nil {
		if strings.Contains(err.Error(), NoRowsInResultSet) {
			return &entities.User{}, domain.ErrUserNotFound
		}

		return &entities.User{}, fmt.Errorf("failed to query user table: %w", err)
	}

	return model.ToEntity(), nil
}

func (u *userDatabase) InsertUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	model := models.NewUserModel(user)

	_, err := u.Client.DB.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), database.DuplicateKeyPrefix) {
			return &entities.User{}, domain.ErrEmailIsAlreadyInUse
		}

		return &entities.User{}, fmt.Errorf("failed to insert into user table: %w", err)
	}

	return model.ToEntity(), nil
}

func (u *userDatabase) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	model := models.NewUserModel(user)

	_, err := u.Client.DB.NewUpdate().Model(model).Set("opt_in = ?", false).Where("id = ?", user.ID).Exec(ctx)
	if err != nil {
		return &entities.User{}, fmt.Errorf("failed to updated user table: %w", err)
	}

	return model.ToEntity(), nil
}

func (u *userDatabase) ListUser(ctx context.Context) ([]*entities.User, error) {
	var modelList []models.User
	query := u.Client.DB.NewSelect().Model(&modelList)

	if err := query.Scan(ctx); err != nil {
		return nil, fmt.Errorf("failed to query user table: %w", err)
	}

	list := make([]*entities.User, 0)
	for _, model := range modelList {
		list = append(list, model.ToEntity())
	}

	return list, nil
}
