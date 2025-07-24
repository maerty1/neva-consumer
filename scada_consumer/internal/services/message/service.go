package message

import (
	"context"

	"scada_consumer/internal/message_broker"
	"scada_consumer/internal/repositories/message"
)

type Service interface {
	RunScadaConsumer(ctx context.Context) error
}

var _ Service = (*service)(nil)

type service struct {
	repository    message.Repository
	messageBroker message_broker.MessageBroker
}

func NewService(repository message.Repository, messageBroker message_broker.MessageBroker) *service {
	return &service{
		repository:    repository,
		messageBroker: messageBroker,
	}
}
