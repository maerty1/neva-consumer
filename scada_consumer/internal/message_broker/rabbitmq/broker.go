package rabbitmq

import (
	"log"
	"scada_consumer/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitmqBroker struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queues  []*amqp.Queue
}

func NewRabbitMQBroker(cfg config.RabbitMQConfig) (*RabbitmqBroker, error) {
	rb := &RabbitmqBroker{}

	err := rb.setupConnection(cfg)
	if err != nil {
		return nil, err
	}

	return rb, nil
}

func (rb *RabbitmqBroker) setupConnection(cfg config.RabbitMQConfig) error {
	conn, err := amqp.Dial(cfg.URL())
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	rb.Conn = conn
	rb.Channel = ch

	for _, queueName := range cfg.ConsumerQueues() {
		q, err := ch.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			return err
		}
		rb.Queues = append(rb.Queues, &q)
	}

	return nil
}

func (rb *RabbitmqBroker) ConsumeMessages() map[string]<-chan amqp.Delivery {
	consumers := make(map[string]<-chan amqp.Delivery)

	for _, q := range rb.Queues {
		msgs, err := rb.Channel.Consume(
			q.Name, // queue
			"",     // consumer
			false,  // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		if err != nil {
			log.Printf("Не удалось зарегистрировать потребителя для очереди %s: %v", q.Name, err)
			continue
		}
		consumers[q.Name] = msgs
	}

	return consumers
}
