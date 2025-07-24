package message

import (
	"context"
	"encoding/json"
	"log"
	"scada_consumer/internal/message_broker/models"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	batchSize    = 100             // Размер батча
	batchTimeout = 5 * time.Second // Таймаут для отправки батча. Автоматическая обработка сообщений, если за batchTimeout не набралось batchSize
)

type QueueType string

// RABBITMQ_ASAP_QUEUE
// RABBITMQ_10MIN_QUEUE
// RABBITMQ_DAILY_QUEUE
const (
	Queue10Min QueueType = "10MIN"
	QueueASAP  QueueType = "ASAP"
	QueueDaily QueueType = "DAILY"
)

func (s *service) RunScadaConsumer(ctx context.Context) error {
	consumers := s.messageBroker.ConsumeMessages()

	var wg sync.WaitGroup

	for queueName, msgs := range consumers {
		wg.Add(1)
		go func(qName string, m <-chan amqp.Delivery) {
			defer wg.Done()
			s.processQueue(ctx, qName, m)
		}(queueName, msgs)
	}

	wg.Wait()
	return nil
}

func (s *service) processQueue(ctx context.Context, queueName string, msgs <-chan amqp.Delivery) {
	var batch []models.Message
	timer := time.NewTimer(batchTimeout)
	defer timer.Stop()

	// Карта для отслеживания последнего времени обработки lastChanged для переменных
	lastProcessed := make(map[string]time.Time)
	var mu sync.Mutex

	for {
		select {
		case <-ctx.Done():
			log.Printf("Контекст отменен, выход из обработчика очереди %s", queueName)
			s.flushBatch(ctx, &batch)
			return
		case msg, ok := <-msgs:
			if !ok {
				// Канал закрыт, обработаем оставшиеся сообщения
				s.flushBatch(ctx, &batch)
				return
			}

			// Скада шлет каждое сообщение в массиве, парсим как массив
			var customMsgs []models.Message
			if err := json.Unmarshal(msg.Body, &customMsgs); err != nil {
				log.Printf("ошибка анмаршаллинга сообщения из очереди %s: %v", queueName, err)
				s.commitMessage(msg)
				continue
			}

			for _, customMsg := range customMsgs {
				if err := customMsg.ParseHash(); err != nil {
					log.Printf("ошибка парсинга хеша сообщения из очереди %s: %v", queueName, err)
					s.commitMessage(msg)
					continue
				}

				// Парсим timestamp из lastChanged
				msgTimestamp, err := time.Parse(time.RFC3339, customMsg.LastChanged)
				if err != nil {
					log.Printf("неверный формат временной метки для сообщения с хешем '%s': %v", customMsg.Hash, err)
					s.commitMessage(msg)
					continue // Пропускаем сообщение с некорректным timestamp
				}

				shouldProcess := false

				mu.Lock()
				lastTime, exists := lastProcessed[customMsg.Variable]
				switch QueueType(queueName) {
				case Queue10Min:
					if !exists || msgTimestamp.After(lastTime.Add(10*time.Minute)) {
						shouldProcess = true
						if msgTimestamp.After(lastTime) {
							lastProcessed[customMsg.Variable] = msgTimestamp
						}
					}
				case QueueDaily:
					if !exists || msgTimestamp.After(lastTime.Add(24*time.Hour)) {
						shouldProcess = true
						if msgTimestamp.After(lastTime) {
							lastProcessed[customMsg.Variable] = msgTimestamp
						}
					}
				case QueueASAP:
					shouldProcess = true
					if msgTimestamp.After(lastTime) {
						lastProcessed[customMsg.Variable] = msgTimestamp
					}
				}
				mu.Unlock()

				if shouldProcess {
					customMsg.RabbitMQMessage = msg
					batch = append(batch, customMsg)

					if len(batch) >= batchSize {
						s.flushBatch(ctx, &batch)
						timer.Reset(batchTimeout)
					}
				} else {
					// Если не обрабатываем, подтверждаем сообщение, но не добавляем в батч
					if commitErr := s.commitMessage(msg); commitErr != nil {
						log.Printf("ошибка коммита сообщения из очереди %s: %v", queueName, commitErr)
					}
				}
			}

		case <-timer.C:
			if len(batch) > 0 {
				s.flushBatch(ctx, &batch)
			}
			timer.Reset(batchTimeout)
		}
	}
}

func (s *service) flushBatch(ctx context.Context, batch *[]models.Message) {
	if len(*batch) == 0 {
		return
	}

	err := s.repository.WriteBatch(ctx, *batch)
	if err != nil {
		log.Printf("ошибка записи батча: %v", err)
		return
	}

	// Подтверждаем все сообщения в батче после успешной записи
	for _, msg := range *batch {
		err := s.commitMessage(msg.RabbitMQMessage)
		if err != nil {
			log.Printf("ошибка коммита сообщения с hash '%s': %v", msg.Hash, err)
		}
	}

	// Переинициализация слайса для освобождения памяти
	*batch = make([]models.Message, 0, batchSize)
}

func (s *service) commitMessage(msg amqp.Delivery) error {
	return msg.Ack(false)
}
