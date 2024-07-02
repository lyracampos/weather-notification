package ports

import (
	"context"
	"weather-notification/internal/domain/entities"
)

type WeatherHTTPGateway interface {
	GetCity(ctx context.Context, city string, state string) (*entities.City, error)
	GetWeather(ctx context.Context, cityID int) (*[]entities.Weather, error)
	GetWeatherCoast(ctx context.Context, cityID int) (*entities.WeatherCoast, error)
}
