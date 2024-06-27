package usecases

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var _ NotifyUsersUseCase = (*notifyUsersUseCase)(nil)

type NotifyUsersInput struct {
	Emails []string
}

type NotifyUsersUseCase interface {
	Execute(ctx context.Context, input NotifyUsersInput) error
}

type notifyUsersUseCase struct {
}

func NewNotifyUsersUseCase() *notifyUsersUseCase {
	return &notifyUsersUseCase{}
}

func (u *notifyUsersUseCase) Execute(ctx context.Context, input NotifyUsersInput) error {
	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@localhost:5672/")
	if err != nil {
		log.Panicf("falied to connect on rabbitmq servers: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"websocket_notification",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("failed to declare a queue: %v", err)
	}

	body := "Hello, World!"
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Fatalf("failed to publish a message: %v", err)
	}

	log.Printf("message sent: %v", body)

	return nil
}
