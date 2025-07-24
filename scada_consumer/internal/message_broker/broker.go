package message_broker

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageBroker интерфейс обеспечивает абстракцию для работы с любым брокером сообщений.
type MessageBroker interface {
	// CommitMessages(ctx context.Context, messages []models.Message) error
	// Это плохо, надо заменить на что-то нейтральное
	ConsumeMessages() map[string]<-chan amqp.Delivery
	// Close() error
}
