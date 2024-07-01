package models

import (
	"time"
	"weather-notification/internal/domain/entities"

	"github.com/uptrace/bun"
)

const UsersTableName = "users"

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64     `bun:"id,pk,autoincrement"`
	FirstName string    `bun:"first_name,notnull"`
	LastName  string    `bun:"last_name,notnull"`
	Email     string    `bun:"email,notnull"`
	Phone     string    `bun:"phone,notnull"`
	CityID    int       `bun:"city_id,notnull"`
	OptIn     bool      `bun:"opt_in,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

func NewUserModel(entity *entities.User) *User {
	return &User{
		BaseModel: bun.BaseModel{},
		ID:        entity.ID,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Email:     entity.Email,
		Phone:     entity.Phone,
		CityID:    entity.CityID,
		OptIn:     entity.OptIn,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func (u *User) ToEntity() *entities.User {
	return &entities.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
		CityID:    u.CityID,
		OptIn:     u.OptIn,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
