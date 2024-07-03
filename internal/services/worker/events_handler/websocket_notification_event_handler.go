package eventshandler

import (
	"context"
	"encoding/json"
	"weather-notification/internal/domain/events"
	"weather-notification/internal/domain/usecases"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var _ WebsocketEventHandler = (*websocketEventHandler)(nil)

type WebsocketEventHandler interface {
	EventHandler(ctx context.Context, msg amqp.Delivery, err error)
}

type websocketEventHandler struct {
	log               *zap.SugaredLogger
	notifyUserUseCase usecases.NotifyUserUseCase
}

func NewWebsocketEventHandler(log *zap.SugaredLogger, notifyUserUseCase usecases.NotifyUserUseCase) websocketEventHandler {
	return websocketEventHandler{
		log:               log,
		notifyUserUseCase: notifyUserUseCase,
	}
}

// nolint: staticcheck
func (h *websocketEventHandler) EventHandler(ctx context.Context, msg amqp.Delivery, err error) {
	event := events.WebsocketNotificationEvent{}
	if err := json.Unmarshal([]byte(msg.Body), &event); err != nil {
		h.log.Errorf("error to unmarshal message to event %v :", err)
	}

	notification, err := h.notifyUserUseCase.Execute(ctx, event.UserEmail)
	if err != nil {
		h.log.Errorf("failed to notify user: %s %v: ", event.UserEmail, err)
	}

	n, err := json.Marshal(notification)
	if err != nil {
		h.log.Errorf("failed to marshal notification user: %v", err)
	}

	h.log.Infof("user %s was notified", string(n))
}
