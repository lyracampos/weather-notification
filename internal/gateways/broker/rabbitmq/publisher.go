package rabbitmq

import (
	"context"
	"fmt"
	"weather-notification/internal/domain/events"
	"weather-notification/internal/domain/ports"
)

var _ ports.PublisherBrokerGateway = (*publisher)(nil)

type publisher struct {
	Client *Client
}

func NewPublisher(client *Client) *publisher {
	return &publisher{
		Client: client,
	}
}

func (n *publisher) WebsocketNotificationEvent(ctx context.Context, event *events.WebsocketNotificationEvent) error {
	if err := n.Client.publish(event, WebsocketNotificatoinQueue); err != nil {
		return fmt.Errorf("failed to enqueue WebsocketNotificationEvent %v: ", err)
	}

	return nil
}

func (n *publisher) EmailNotificationEvent(ctx context.Context, event *events.EmailNotificationEvent) error {
	if err := n.Client.publish(event, EmailNotificatoinQueue); err != nil {
		return fmt.Errorf("failed to enqueue EmailNotificationEvent %v: ", err)
	}

	return nil
}

func (n *publisher) SMSNotificationEvent(ctx context.Context, event *events.SMSNotificationEvent) error {
	if err := n.Client.publish(event, SMSNotificatoinQueue); err != nil {
		return fmt.Errorf("failed to enqueue SMSNotificationEvent %v: ", err)
	}

	return nil
}

func (n *publisher) PushNotificationEvent(ctx context.Context, event *events.PushNotificationEvent) error {
	if err := n.Client.publish(event, PushNotificatoinQueue); err != nil {
		return fmt.Errorf("failed to enqueue PushNotificationEvent %v: ", err)
	}

	return nil
}
