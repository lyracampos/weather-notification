package rabbitmq

import (
	"encoding/json"
	"fmt"
	"weather-notification/configs"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	WebsocketNotificatoinQueue = "websocket-notifications"
	EmailNotificatoinQueue     = "email-notifications"
	SMSNotificatoinQueue       = "sms-notifications"
	PushNotificatoinQueue      = "push-notifications"
)

type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewClient(log *zap.SugaredLogger, config *configs.Config) (*Client, error) {
	client := &Client{}
	var err error

	client.conn, err = amqp.Dial(config.Broker.ConnectionURL)
	if err != nil {
		return &Client{}, fmt.Errorf("failed to connect on rabbitmq server: %w", err)
	}

	client.ch, err = client.conn.Channel()
	if err != nil {
		return &Client{}, fmt.Errorf("failed to create rabbitmq channel: %w", err)
	}
	err = client.configureQueues()
	if err != nil {
		return &Client{}, fmt.Errorf("failed to create rabbitmq queues: %w", err)
	}

	log.Info("rabbitmq - client started...")

	return client, nil
}

func (c *Client) Close() {
	c.ch.Close()
	c.conn.Close()
}

func (c *Client) configureQueues() error {
	var err error
	if err = c.configureQueue(WebsocketNotificatoinQueue); err != nil {
		return fmt.Errorf("failed to create websocket-notifications queue: %w", err)
	}
	if err = c.configureQueue(EmailNotificatoinQueue); err != nil {
		return fmt.Errorf("failed to create email-notifications queue: %w", err)
	}
	if err = c.configureQueue(SMSNotificatoinQueue); err != nil {
		return fmt.Errorf("failed to create sms-notifications queue: %w", err)
	}
	if err = c.configureQueue(PushNotificatoinQueue); err != nil {
		return fmt.Errorf("failed to create push-notifications queue: %w", err)
	}

	return nil
}

func (c *Client) configureQueue(queueName string) error {
	_, err := c.ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return err
}

func (c *Client) publish(event interface{}, queue string) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to convert event struct to json: %w", err)
	}

	_ = c.ch.Publish("", queue, true, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(eventJSON),
	})

	return nil
}
