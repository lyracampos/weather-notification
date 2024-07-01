package eventshandler

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var _ WebsocketEventHandler = (*websocketEventHandler)(nil)

type WebsocketEventHandler interface {
	EventHandler(msg amqp.Delivery, err error)
}

type websocketEventHandler struct {
	log *zap.SugaredLogger
}

func NewWebsocketEventHandler(log *zap.SugaredLogger) websocketEventHandler {
	return websocketEventHandler{
		log: log,
	}
}

func (u *websocketEventHandler) EventHandler(msg amqp.Delivery, err error) {
	u.log.Info("webscoket event handler")
	u.log.Info("message received %s", string(msg.Body))
}
