package usecases

import (
	"context"
	"fmt"
	"weather-notification/internal/domain/entities"
	"weather-notification/internal/domain/ports"

	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

var _ EnqueueNotificationsUseCase = (*enqueueNotificationsUseCase)(nil)

type QueueNotificationsInput struct {
	Users []string `validate:"required,min=1"`
	Type  string   `validate:"required,oneof=web push email sms"`
}

func (u *QueueNotificationsInput) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

type EnqueueNotificationsUseCase interface {
	Execute(ctx context.Context, input QueueNotificationsInput) error
}

type enqueueNotificationsUseCase struct {
	log                *zap.SugaredLogger
	notificationBroker ports.NotificationBrokerGateway
}

func NewEnqueueNotificationsUseCase(log *zap.SugaredLogger, notificationBroker ports.NotificationBrokerGateway) *enqueueNotificationsUseCase {
	return &enqueueNotificationsUseCase{
		log:                log,
		notificationBroker: notificationBroker,
	}
}

func (u *enqueueNotificationsUseCase) Execute(ctx context.Context, input QueueNotificationsInput) error {
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	for _, user := range input.Users {
		go u.enqueueNotification(ctx, user, input.Type)
	}

	// errCh := make(chan error, len(input.Users))

	// for _, user := range input.Users {
	// 	go u.enqueueNotification(ctx, user, input.Type, errCh)
	// }

	return nil
}

func (u *enqueueNotificationsUseCase) enqueueNotification(ctx context.Context, emailTo string, notificationType string) error {
	notification := entities.NewNotification(emailTo, entities.NotificationTypeWebSocket, entities.NotificationStatusQueued)

	if notificationType == "web" {
		err := u.notificationBroker.PublishWebNotification(ctx, notification)
		if err != nil {
			u.log.Errorf("failed to enqueue websocket notification to user: %s :%w", emailTo, err)
		}
	}

	return nil
}

// func (u *queueNotificationsUseCase) enqueueNotification(ctx context.Context, emailTo string, notificationType string, errCh chan<- error) {
// 	notification := entities.NewNotification(emailTo, entities.NotificationTypeWebSocket, entities.NotificationStatusQueued)
// 	err := notification.Validate()
// 	if err != nil {
// 		errCh <- fmt.Errorf("invalid request: %w", err)
// 	}

// 	if notificationType == "web" {
// 		err = u.notificationBroker.PublishWebNotification(ctx, notification)
// 		if err != nil {
// 			errCh <- fmt.Errorf("failed to enqueue websocket notification: %w", err)
// 		}
// 	}
// }