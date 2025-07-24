package app

import (
	"scada_consumer/internal/db"
	"scada_consumer/internal/message_broker"
	messageRepository "scada_consumer/internal/repositories/message"
	messageService "scada_consumer/internal/services/message"
)

type serviceProvider struct {
	postgresDB        db.PostgresClient
	rabbitBroker      message_broker.MessageBroker
	messageRepository messageRepository.Repository
	messageService    messageService.Service
}

func newServiceProvider(db db.PostgresClient, rabbitBroker message_broker.MessageBroker) *serviceProvider {
	return &serviceProvider{
		postgresDB:   db,
		rabbitBroker: rabbitBroker,
	}
}

func (s *serviceProvider) PostgresDB() db.PostgresClient {
	return s.postgresDB
}

func (s *serviceProvider) RabbitBroker() message_broker.MessageBroker {
	return s.rabbitBroker
}

func (s *serviceProvider) MessageRepository() messageRepository.Repository {
	if s.messageRepository == nil {
		s.messageRepository = messageRepository.NewRepository(s.PostgresDB())
	}
	return s.messageRepository
}

func (s *serviceProvider) MessageService() messageService.Service {
	if s.messageService == nil {
		s.messageService = messageService.NewService(s.MessageRepository(), s.RabbitBroker())
	}
	return s.messageService
}
