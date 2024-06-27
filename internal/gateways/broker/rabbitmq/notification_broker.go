package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"weather-notification/internal/domain/entities"
	"weather-notification/internal/domain/ports"

	amqp "github.com/rabbitmq/amqp091-go"
)

var _ ports.NotificationBrokerGateway = (*notificationBroker)(nil)

type notificationBroker struct {
	Client *Client
}

func NewNotificationBroker(client *Client) *notificationBroker {
	return &notificationBroker{
		Client: client,
	}
}

func (n *notificationBroker) PublishWebNotification(ctx context.Context, notification *entities.Notification) error {
	messageJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to convert notification struct to json: %w", err)
	}

	_ = n.Client.ch.Publish("", "websocket-notifications", true, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(messageJSON),
	})

	return nil
}
