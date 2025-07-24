package config

import (
	"errors"
	"os"
	"strings"
)

const (
	kafkaBrokersEnvName          = "KAFKA_BROKERS"
	kafkaConsumerTopicName       = "KAFKA_CONSUMER_TOPIC"
	kafkaConsumerGroupIDEnvName  = "KAFKA_CONSUMER_GROUP_ID"
	kafkaSecurityProtocolEnvName = "KAFKA_SECURITY_PROTOCOL"
	kafkaSaslMechanismEnvName    = "KAFKA_SASL_MECHANISM"
	kafkaSaslUsernameEnvName     = "KAFKA_SASL_PLAIN_USERNAME"
	kafkaSaslPasswordEnvName     = "KAFKA_SASL_PLAIN_PASSWORD"
)

type KafkaConfig interface {
	Brokers() []string
	ConsumerGroupID() string
	SecurityProtocol() string
	SaslMechanism() string
	SaslUsername() string
	SaslPassword() string
	ConsumerTopic() string
}

type kafkaConfig struct {
	brokers          []string
	consumerTopic    string
	consumerGroupID  string
	securityProtocol string
	saslMechanism    string
	saslUsername     string
	saslPassword     string
}

// GetKafkaConfig создает и возвращает конфигурацию Kafka на основе переменных среды.
func GetKafkaConfig() (KafkaConfig, error) {
	brokers := os.Getenv(kafkaBrokersEnvName)
	if brokers == "" {
		return nil, errors.New("адреса брокеров Kafka не указаны")
	}

	return &kafkaConfig{
		brokers:          strings.Split(brokers, ","),
		consumerTopic:    os.Getenv(kafkaConsumerTopicName),
		consumerGroupID:  os.Getenv(kafkaConsumerGroupIDEnvName),
		securityProtocol: os.Getenv(kafkaSecurityProtocolEnvName),
		saslMechanism:    os.Getenv(kafkaSaslMechanismEnvName),
		saslUsername:     os.Getenv(kafkaSaslUsernameEnvName),
		saslPassword:     os.Getenv(kafkaSaslPasswordEnvName),
	}, nil
}

func (cfg *kafkaConfig) Brokers() []string {
	return cfg.brokers
}

func (cfg *kafkaConfig) ConsumerGroupID() string {
	return cfg.consumerGroupID
}

func (cfg *kafkaConfig) SecurityProtocol() string {
	return cfg.securityProtocol
}

func (cfg *kafkaConfig) SaslMechanism() string {
	return cfg.saslMechanism
}

func (cfg *kafkaConfig) SaslUsername() string {
	return cfg.saslUsername
}

func (cfg *kafkaConfig) SaslPassword() string {
	return cfg.saslPassword
}

func (cfg *kafkaConfig) ConsumerTopic() string {
	return cfg.consumerTopic
}

func NewKafkaConfig(brokers []string, consumerTopic, consumerGroupID, securityProtocol, saslMechanism, saslUsername, saslPassword string) KafkaConfig {
	return &kafkaConfig{
		brokers:          brokers,
		consumerTopic:    consumerTopic,
		consumerGroupID:  consumerGroupID,
		securityProtocol: securityProtocol,
		saslMechanism:    saslMechanism,
		saslUsername:     saslUsername,
		saslPassword:     saslPassword,
	}
}
