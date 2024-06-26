package postgres

import (
	"database/sql"
	"fmt"
	"weather-notification/configs"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.uber.org/zap"
)

// const (
// 	DuplicateKeyPrefix = "duplicate key value violates unique constraint"
// 	NoRowsInResultSet  = "no rows in result set"
// )

type Client struct {
	log    *zap.SugaredLogger
	DB     *bun.DB
	config *configs.Config
}

func NewClient(log *zap.SugaredLogger, config *configs.Config) (*Client, error) {
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.UsersAPI.Database.ConnectionString)))

	newDB := bun.NewDB(sqlDB, pgdialect.New())
	// newDB.AddQueryHook(bunotel.NewQueryHook())

	if err := newDB.Ping(); err != nil {
		return &Client{}, fmt.Errorf("failed to connect to postgres database: %w", err)
	}

	log.Info("postgres client started successfully")

	return &Client{
		DB:     newDB,
		config: config,
		log:    log,
	}, nil
}
