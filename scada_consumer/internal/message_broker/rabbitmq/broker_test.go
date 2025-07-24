// rabbitmq_test/rabbitmq_test.go

package rabbitmq_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"scada_consumer/internal/app"
	"scada_consumer/internal/message_broker/models"
	"scada_consumer/internal/message_broker/rabbitmq"
	"scada_consumer/tests"
	"testing"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
)

var brokerTest *rabbitmq.RabbitmqBroker
var appInstance *app.App

// TestMain выполняется перед запуском всех тестов.
// Здесь происходит инициализация приложения и брокера RabbitMQ.
func TestMain(m *testing.M) {
	fmt.Println("Start tests...")
	ctx := context.Background()

	tests.Init(ctx)

	app, cleanup, err := tests.GetApp()
	if err != nil {
		log.Fatalf("Ошибка получения app instance: %s", err)
	}

	appInstance = app
	brokerTest = app.S().RabbitBroker().(*rabbitmq.RabbitmqBroker)

	exitCode := m.Run()

	cleanup()

	os.Exit(exitCode)
}

// setupTestQueue объявляет новую очередь для теста.
func setupTestQueue(t *testing.T, queueName string) {
	_, err := brokerTest.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	require.NoError(t, err, "Не удалось объявить очередь %s: %v", queueName, err)
}

// cleanupTestQueue удаляет очередь после завершения теста.
func cleanupTestQueue(queueName string) {
	_, err := brokerTest.Channel.QueueDelete(queueName, false, false, false)
	if err != nil {
		log.Printf("Ошибка удаления очереди %s: %v", queueName, err)
	}
}

// publishMessages публикует набор сообщений в указанную очередь.
func publishMessages(t *testing.T, queueName string, messages []models.Message) {
	for _, message := range messages {
		messageBytes, err := json.Marshal(message)
		require.NoError(t, err, "Ошибка маршаллинга сообщения: %v", err)

		err = brokerTest.Channel.Publish(
			"",        // exchange
			queueName, // routing key (queue name)
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        messageBytes,
			},
		)
		require.NoError(t, err, "Ошибка публикации сообщения в очередь %s: %v", queueName, err)
	}
}

// receiveAndAcknowledgeMessages принимает очередь и ожидает получения определённого количества сообщений.
// После успешной обработки сообщения оно подтверждается.
func receiveAndAcknowledgeMessages(t *testing.T, queueName string, expectedMessages []models.Message) {
	// Создаём отдельный канал для потребления сообщений из данной очереди
	ch, err := brokerTest.Conn.Channel()
	require.NoError(t, err, "Не удалось открыть канал для потребления сообщений")

	defer ch.Close()

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	require.NoError(t, err, "Не удалось зарегистрировать потребителя для очереди %s: %v", queueName, err)

	receivedMessages := make([]models.Message, 0, len(expectedMessages))

	timeout := time.After(10 * time.Second)
	for i := 0; i < len(expectedMessages); i++ {
		select {
		case msg, ok := <-msgs:
			if !ok {
				t.Fatalf("Канал для очереди %s закрыт до получения всех сообщений", queueName)
			}

			var received models.Message
			err := json.Unmarshal(msg.Body, &received)
			require.NoError(t, err, "Ошибка анмаршаллинга сообщения из очереди %s: %v", queueName, err)

			receivedMessages = append(receivedMessages, received)

			// Подтверждаем получение сообщения
			err = msg.Ack(false)
			require.NoError(t, err, "Ошибка подтверждения сообщения из очереди %s: %v", queueName, err)

		case <-timeout:
			t.Fatal("Timeout: Не удалось получить все сообщения из очереди")
		}
	}

	// Проверяем, что все ожидаемые сообщения были получены
	require.Len(t, receivedMessages, len(expectedMessages), "Ожидалось получить %d сообщений, получено %d", len(expectedMessages), len(receivedMessages))
	for i, msg := range receivedMessages {
		require.Equal(t, expectedMessages[i].Hash, msg.Hash, "Hash сообщения не совпадает")
	}
}

// TestRabbitMQBroker_Connection проверяет, что соединение и канал с RabbitMQ установлены корректно.
func TestRabbitMQBroker_Connection(t *testing.T) {
	require.NotNil(t, brokerTest, "Broker должен быть инициализирован")
	require.NotNil(t, brokerTest.Conn, "Соединение должно быть установлено")
	require.NotNil(t, brokerTest.Channel, "Канал должен быть открыт")
}

// TestRabbitMQBroker_ConsumeMessages проверяет, что сообщения могут быть потреблены из очереди и подтверждены.
func TestRabbitMQBroker_ConsumeMessages(t *testing.T) {
	queueName := "test_queue_consume_" + uuid.New().String()
	setupTestQueue(t, queueName)
	defer cleanupTestQueue(queueName)

	messageData := models.Message{
		Value:              "True",
		DataType:           "DINT",
		LastChanged:        "2024-04-04T07:04:15.5548405Z",
		StatusCodes:        12,
		NodeId:             "f80f7085-1dbf-4f49-9e63-a279a4a3227e",
		NodeName:           "RabbitMQ super",
		OwnerId:            1,
		Hash:               "Hash",
		DataPointClassEnum: "Input",
	}

	publishMessages(t, queueName, []models.Message{messageData})
	receiveAndAcknowledgeMessages(t, queueName, []models.Message{messageData})
}

// TestRabbitMQBroker_ConsumeMultipleMessages проверяет обработку нескольких сообщений.
func TestRabbitMQBroker_ConsumeMultipleMessages(t *testing.T) {
	queueName := "test_queue_multiple_" + uuid.New().String()
	setupTestQueue(t, queueName)
	defer cleanupTestQueue(queueName)

	messages := []models.Message{
		{
			Value:              "True",
			DataType:           "DINT",
			LastChanged:        "2024-04-04T07:04:15.5548405Z",
			StatusCodes:        12,
			NodeId:             "f80f7085-1dbf-4f49-9e63-a279a4a3227e",
			NodeName:           "RabbitMQ super",
			OwnerId:            1,
			Hash:               "Hash1",
			DataPointClassEnum: "Input",
		},
		{
			Value:              "False",
			DataType:           "BOOL",
			LastChanged:        "2024-04-04T07:05:16.5548405Z",
			StatusCodes:        0,
			NodeId:             "a1234567-1dbf-4f49-9e63-a279a4a3227f",
			NodeName:           "RabbitMQ test",
			OwnerId:            2,
			Hash:               "Hash2",
			DataPointClassEnum: "Output",
		},
		{
			Value:              "42",
			DataType:           "INT",
			LastChanged:        "2024-04-04T07:06:17.5548405Z",
			StatusCodes:        99,
			NodeId:             "b2345678-1dbf-4f49-9e63-a279a4a3227g",
			NodeName:           "RabbitMQ another",
			OwnerId:            3,
			Hash:               "Hash3",
			DataPointClassEnum: "Intermediate",
		},
	}

	publishMessages(t, queueName, messages)
	receiveAndAcknowledgeMessages(t, queueName, messages)
}

// TestRabbitMQBroker_NoDuplicateAfterAck проверяет, что сообщения не дублируются после подтверждения.
func TestRabbitMQBroker_NoDuplicateAfterAck(t *testing.T) {
	queueName := "test_queue_no_duplicate_" + uuid.New().String()
	setupTestQueue(t, queueName)
	defer cleanupTestQueue(queueName)

	messageData := models.Message{
		Value:              "UniqueValue",
		DataType:           "STRING",
		LastChanged:        "2024-04-04T08:00:00Z",
		StatusCodes:        1,
		NodeId:             "c3456789-1dbf-4f49-9e63-a279a4a3227h",
		NodeName:           "RabbitMQ unique",
		OwnerId:            4,
		Hash:               "UniqueHash",
		DataPointClassEnum: "Unique",
	}

	publishMessages(t, queueName, []models.Message{messageData})
	receiveAndAcknowledgeMessages(t, queueName, []models.Message{messageData})

	// Пытаемся получить сообщение снова, ожидаем, что его нет
	ch, err := brokerTest.Conn.Channel()
	require.NoError(t, err, "Не удалось открыть канал для проверки дублирования сообщений")
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	require.NoError(t, err, "Не удалось зарегистрировать потребителя для очереди %s: %v", queueName, err)

	select {
	case msg, ok := <-msgs:
		if ok {
			var received models.Message
			err := json.Unmarshal(msg.Body, &received)
			require.NoError(t, err, "Ошибка анмаршаллинга сообщения из очереди %s: %v", queueName, err)
			require.Equal(t, messageData.Hash, received.Hash, "Получено дублирующее сообщение с hash: %s", received.Hash)
			msg.Ack(false)
			t.Fatalf("Получено дублирующее сообщение: %s", msg.Body)
		}
	case <-time.After(2 * time.Second):
		// Успешно: сообщения не было
	}
}
