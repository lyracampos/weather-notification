package ports

import (
	"context"
	"weather-notification/internal/domain/entities"
)

type NotificationBrokerGateway interface {
	PublishWebNotification(ctx context.Context, notification *entities.Notification) error
}

type NotificationDatabaseGateway interface {
	InsertNotification(ctx context.Context, notification *entities.Notification) (*entities.Notification, error)
}
