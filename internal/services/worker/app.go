package worker

import (
	"context"
	"log"
	"weather-notification/configs"
	"weather-notification/internal/domain/usecases"
	"weather-notification/internal/gateways/broker/rabbitmq"
	"weather-notification/internal/gateways/database/postgres"
	api "weather-notification/internal/gateways/http"
	eventshandler "weather-notification/internal/services/worker/events_handler"

	"go.uber.org/zap"
)

func Run(config *configs.Config, workerType string) {
	ctx := context.Background()
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatalf("failed to defer logger sync: %v", err)
		}
	}()

	sugar := logger.Sugar()

	// initialize client dependences
	databaseClient, err := postgres.NewClient(sugar, config)
	if err != nil {
		sugar.Fatalf("failed to initialize postgres client: %w", err)
	}
	defer databaseClient.Close()
	userDatabase := postgres.NewUserDatabase(databaseClient)

	brokerClient, err := rabbitmq.NewClient(sugar, config)
	if err != nil {
		sugar.Fatalf("failed to initialize rabbitmq client: %w", err)
	}
	defer brokerClient.Close()

	weatherAPI := api.NewWeatherAPI(sugar, config)
	webNotificationAPI := api.NewWebSocketServer(sugar, config)

	notifyUserUseCase := usecases.NewNotifyUserUseCase(sugar, userDatabase, weatherAPI, webNotificationAPI)

	switch workerType {
	case "web":
		runWebNotificationWorker(ctx, sugar, brokerClient, notifyUserUseCase)
	default:
		sugar.Error("invalid worker type")
	}
}

func runWebNotificationWorker(ctx context.Context, log *zap.SugaredLogger, brokerClient *rabbitmq.Client, notifyUserUseCase usecases.NotifyUserUseCase) {
	log.Info("running worker for websocket notifications")
	websocketEventHandler := eventshandler.NewWebsocketEventHandler(log, notifyUserUseCase)
	consumerWebsocket := rabbitmq.NewConsumerWebsocket(log, brokerClient, websocketEventHandler.EventHandler)
	consumerWebsocket.Consume(ctx)
}
