package config

import (
	"errors"
	"os"
)

const (
	rabbitmqURLEnvName           = "RABBITMQ_URL"
	rabbitmqConsumerQueueEnvName = "RABBITMQ_CONSUMER_QUEUE"
	rabbitmqConsumerTagEnvName   = "RABBITMQ_CONSUMER_TAG"
	rabbitmqAsapQueueEnvName     = "RABBITMQ_ASAP_QUEUE"
	rabbitmq10MinQueueEnvName    = "RABBITMQ_10MIN_QUEUE"
	rabbitmqDailyQueueEnvName    = "RABBITMQ_DAILY_QUEUE"
)

type RabbitMQConfig interface {
	URL() string
	// ConsumerQueue() string
	ConsumerTag() string
	ConsumerQueues() []string
	ConsumerQueue(index int) string
}

type rabbitmqConfig struct {
	url            string
	consumerQueue  string
	consumerTag    string
	consumerQueues []string
}

func GetRabbitMQConfig() (RabbitMQConfig, error) {
	url := os.Getenv(rabbitmqURLEnvName)
	if url == "" {
		return nil, errors.New("URL RabbitMQ не указан")
	}

	consumerQueue := os.Getenv(rabbitmqConsumerQueueEnvName)
	if consumerQueue == "" {
		return nil, errors.New("очередь RabbitMQ не указана")
	}

	queueAsap := os.Getenv(rabbitmqAsapQueueEnvName)
	if queueAsap == "" {
		return nil, errors.New("queueAsap не указан")
	}
	queue10Min := os.Getenv(rabbitmq10MinQueueEnvName)
	if queue10Min == "" {
		return nil, errors.New("queue10Min не указан")
	}
	queueDaily := os.Getenv(rabbitmqDailyQueueEnvName)
	if queueDaily == "" {
		return nil, errors.New("queueDaily не указан")
	}

	consumerQueues := []string{queue10Min, queueAsap, queueDaily}
	return &rabbitmqConfig{
		url:            url,
		consumerQueue:  consumerQueue,
		consumerTag:    os.Getenv(rabbitmqConsumerTagEnvName),
		consumerQueues: consumerQueues,
	}, nil
}

func (cfg *rabbitmqConfig) URL() string {
	return cfg.url
}

func (cfg *rabbitmqConfig) ConsumerTag() string {
	return cfg.consumerTag
}

func (cfg *rabbitmqConfig) ConsumerQueues() []string {
	return cfg.consumerQueues
}

func (cfg *rabbitmqConfig) ConsumerQueue(index int) string {
	if index < len(cfg.consumerQueues) {
		return cfg.consumerQueues[index]
	}
	return ""
}

func NewRabbitMQConfig(url, consumerQueue, consumerTag string, consumerQueues []string) RabbitMQConfig {
	return &rabbitmqConfig{
		url:            url,
		consumerQueue:  consumerQueue,
		consumerTag:    consumerTag,
		consumerQueues: consumerQueues,
	}
}
