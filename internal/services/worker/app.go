package worker

import (
	"context"
	"fmt"
	"log"
	"weather-notification/configs"
	"weather-notification/internal/gateways/broker/rabbitmq"
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

	brokerClient, err := rabbitmq.NewClient(sugar, config)
	if err != nil {
		sugar.Fatalf("failed to initialize rabbitmq client: %w", err)
	}
	defer brokerClient.Close()

	switch workerType {
	case "websocket":
		runWebsocketNotificationWorker(ctx, sugar, brokerClient)
	case "email":
		runEmailNotificationWorker(config)
	case "sms":
		runSMSNotificationWorker(config)
	case "push":
		runPushNotificationWorker(config)
	}
}

func runWebsocketNotificationWorker(ctx context.Context, log *zap.SugaredLogger, brokerClient *rabbitmq.Client) {
	log.Info("running websocket notification worker")
	websocketEventHandler := eventshandler.NewWebsocketEventHandler(log)
	consumerWebsocket := rabbitmq.NewConsumerWebsocket(log, brokerClient, websocketEventHandler.EventHandler)
	consumerWebsocket.Consume(ctx)
	// consumerWebsocket := rabbitmq.NewConsumerWebsocket(log, brokerClient, websocketEventHandler)

	// websocketConsumer := consumers.NewWebsocketConsumer(log, consumerBroker)
	// websocketConsumer.Consume(ctx)
}

func runEmailNotificationWorker(config *configs.Config) {
	fmt.Printf("running email notification worker")
}

func runSMSNotificationWorker(config *configs.Config) {
	fmt.Printf("running sms notification worker")
}

func runPushNotificationWorker(config *configs.Config) {
	fmt.Printf("running push notification worker")
}
