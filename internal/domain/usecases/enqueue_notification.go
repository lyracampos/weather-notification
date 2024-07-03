package usecases

import (
	"context"
	"fmt"
	"weather-notification/internal/domain/entities"
	"weather-notification/internal/domain/events"
	"weather-notification/internal/domain/ports"

	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

var _ EnqueueNotificationsUseCase = (*enqueueNotificationsUseCase)(nil)

type EnqueueNotificationsInput struct {
	Users []string `validate:"required,min=1"`
	Type  string   `validate:"required,oneof=websocket push email sms"`
}

func (u *EnqueueNotificationsInput) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

type EnqueueNotificationsUseCase interface {
	Execute(ctx context.Context, input EnqueueNotificationsInput) error
}

type enqueueNotificationsUseCase struct {
	log           *zap.SugaredLogger
	publishBroker ports.PublisherBrokerGateway
}

func NewEnqueueNotificationsUseCase(log *zap.SugaredLogger, publishBroker ports.PublisherBrokerGateway) *enqueueNotificationsUseCase {
	return &enqueueNotificationsUseCase{
		log:           log,
		publishBroker: publishBroker,
	}
}

// nolint: errcheck
func (u *enqueueNotificationsUseCase) Execute(ctx context.Context, input EnqueueNotificationsInput) error {
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	for _, user := range input.Users {
		go u.enqueueNotification(ctx, user, input.Type)
	}

	return nil
}

func (u *enqueueNotificationsUseCase) enqueueNotification(ctx context.Context, emailTo string, notificationType string) error {
	switch notificationType {
	case entities.NotificationTypeWebSocket:
		event := events.WebsocketNotificationEvent{UserEmail: emailTo}
		if err := u.publishBroker.WebsocketNotificationEvent(ctx, &event); err != nil {
			u.log.Errorf("failed to enqueue websockt notification event for user: %s %v", event.UserEmail, err)
		}
	case entities.NotificationTypeEmail:
		event := events.EmailNotificationEvent{UserEmail: emailTo}
		if err := u.publishBroker.EmailNotificationEvent(ctx, &event); err != nil {
			u.log.Errorf("failed to enqueue email notification event for user: %s %v", event.UserEmail, err)
		}
	case entities.NotificationTypeSMS:
		event := events.SMSNotificationEvent{UserEmail: emailTo}
		if err := u.publishBroker.SMSNotificationEvent(ctx, &event); err != nil {
			u.log.Errorf("failed to enqueue sms notification event for user: %s %v", event.UserEmail, err)
		}
	case entities.NotificationTypePush:
		event := events.PushNotificationEvent{UserEmail: emailTo}
		if err := u.publishBroker.PushNotificationEvent(ctx, &event); err != nil {
			u.log.Errorf("failed to enqueue push notification event for user: %s %v", event.UserEmail, err)
		}
	default:
		return nil
	}

	return nil
}
