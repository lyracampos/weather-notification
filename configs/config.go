package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Config struct {
		UsersAPI UsersAPI
	}

	UsersAPI struct {
		API      API
		Database Database
	}
	API struct {
		Host         string
		Port         int
		WriteTimeout int
		ReadTimeout  int
		IdleTimeout  int
	}
	Database struct {
		ConnectionString string
	}
)

func NewConfig(configFilePath string) (*Config, error) {
	config, err := readConfig(configFilePath)
	if err != nil {
		return &Config{}, fmt.Errorf("failed to get viper config: %w", err)
	}

	return &Config{
		UsersAPI: UsersAPI{
			API: API{
				Host:         config.GetString("users-api.api.host"),
				Port:         config.GetInt("users-api.api.port"),
				WriteTimeout: config.GetInt("users-api.api.writeTimeout"),
				ReadTimeout:  config.GetInt("users-api.api.readTimeout"),
				IdleTimeout:  config.GetInt("users-api.api.idleTimeout"),
			},
			Database: Database{
				ConnectionString: config.GetString("users-api.database.connectionString"),
			},
		},
	}, nil
}

func readConfig(configFilePath string) (*viper.Viper, error) {
	config := viper.New()
	config.SetConfigType("yaml")
	config.SetConfigFile(configFilePath)
	config.AutomaticEnv()

	if err := config.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return config, nil
}
