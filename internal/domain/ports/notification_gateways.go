package ports

import (
	"context"
	"weather-notification/internal/domain/entities"
	"weather-notification/internal/domain/events"
)

type PublisherBrokerGateway interface {
	WebsocketNotificationEvent(ctx context.Context, event *events.WebsocketNotificationEvent) error
	EmailNotificationEvent(ctx context.Context, event *events.EmailNotificationEvent) error
	SMSNotificationEvent(ctx context.Context, event *events.SMSNotificationEvent) error
	PushNotificationEvent(ctx context.Context, event *events.PushNotificationEvent) error
}

type ConsumerWebsocketGateway interface {
	Consume(ctx context.Context)
}

type NotificationDatabaseGateway interface {
	InsertNotification(ctx context.Context, notification *entities.Notification) (*entities.Notification, error)
}

type WebNotificationHTTPGateway interface {
	SendNotification(ctx context.Context, user *entities.User, weather *[]entities.Weather, weatherCoast *entities.WeatherCoast) error
}
