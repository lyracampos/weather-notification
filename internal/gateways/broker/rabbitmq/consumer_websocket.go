package rabbitmq

import (
	"context"
	"weather-notification/internal/domain/ports"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var _ ports.ConsumerWebsocketGateway = (*consumerWebsocket)(nil)

type consumerWebsocket struct {
	log          *zap.SugaredLogger
	Client       *Client
	eventHandler func(msg amqp.Delivery, err error)
}

func NewConsumerWebsocket(log *zap.SugaredLogger, client *Client, eventHandler func(msg amqp.Delivery, err error)) *consumerWebsocket {
	return &consumerWebsocket{
		log:          log,
		Client:       client,
		eventHandler: eventHandler,
	}
}

func (c *consumerWebsocket) OnError(err error, msg string) {
	if err != nil {
		c.eventHandler(amqp.Delivery{}, err)
	}
}

func (c *consumerWebsocket) Consume(ctx context.Context) {
	msgs, err := c.Client.ch.Consume(WebsocketNotificatoinQueue, "", true, false, false, false, nil)
	c.OnError(err, "failed to register a websocket consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			c.eventHandler(d, nil)
		}
	}()
	c.log.Info("Started listening for messages on websocket-notifications queue")
	<-forever
}
