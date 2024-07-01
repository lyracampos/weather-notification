package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Config struct {
		API        API
		Database   Database
		Broker     Broker
		WeatherAPI WeatherAPI
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
	Broker struct {
		ConnectionURL string
	}
	WeatherAPI struct {
		AddressURL string
		Timeout    int
		Retries    int
	}
)

func NewConfig(configFilePath string) (*Config, error) {
	config, err := readConfig(configFilePath)
	if err != nil {
		return &Config{}, fmt.Errorf("failed to get viper config: %w", err)
	}

	return &Config{
		API: API{
			Host:         config.GetString("api.host"),
			Port:         config.GetInt("api.port"),
			WriteTimeout: config.GetInt("api.writeTimeout"),
			ReadTimeout:  config.GetInt("api.readTimeout"),
			IdleTimeout:  config.GetInt("api.idleTimeout"),
		},
		Database: Database{
			ConnectionString: config.GetString("database.connectionString"),
		},
		Broker: Broker{
			ConnectionURL: config.GetString("broker.connectionURL"),
		},
		WeatherAPI: WeatherAPI{
			AddressURL: config.GetString("weatherAPI.addressURL"),
			Timeout:    config.GetInt("weatherAPI.timeout"),
			Retries:    config.GetInt("weatherAPI.retries"),
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
