package usecases

import (
	"context"
	"fmt"
	"weather-notification/internal/domain/entities"
	"weather-notification/internal/domain/ports"

	"go.uber.org/zap"
)

var _ NotifyUserUseCase = (*notifyUserUseCase)(nil)

type NotifyUserUseCaseOutput struct {
	User         *entities.User
	Weather      *[]entities.Weather
	WeatherCoast *entities.WeatherCoast
}

type NotifyUserUseCase interface {
	Execute(ctx context.Context, email string) (*NotifyUserUseCaseOutput, error)
}

type notifyUserUseCase struct {
	log          *zap.SugaredLogger
	userDatabase ports.UserDatabaseGateway
	weatherAPI   ports.WeatherHTTPGateway
}

func NewNotifyUserUseCase(log *zap.SugaredLogger, userDatabase ports.UserDatabaseGateway, weatherAPI ports.WeatherHTTPGateway) *notifyUserUseCase {
	return &notifyUserUseCase{
		log:          log,
		userDatabase: userDatabase,
		weatherAPI:   weatherAPI,
	}
}

func (u *notifyUserUseCase) Execute(ctx context.Context, email string) (*NotifyUserUseCaseOutput, error) {
	user, err := u.userDatabase.GetUser(ctx, email)
	if err != nil {
		return &NotifyUserUseCaseOutput{}, fmt.Errorf("failed to get user from database")
	}

	weather, err := u.weatherAPI.GetWeather(ctx, user.CityID)
	if err != nil {
		u.log.Infof("failed to get weather from weather API :%v", err)
	}

	weatherCoast, err := u.weatherAPI.GetWeatherCoast(ctx, user.CityID)
	if err != nil {
		u.log.Infof("failed to get weather coast from weather API :%v", err)
	}

	return &NotifyUserUseCaseOutput{
		User:         user,
		Weather:      weather,
		WeatherCoast: weatherCoast,
	}, nil
}
