package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"scada_consumer/internal/app"
	"scada_consumer/internal/config"
	"scada_consumer/internal/db"
	"scada_consumer/internal/message_broker/rabbitmq"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработка сигналов ОС
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		cancel() // Отправляем сигнал для отмены контекста
	}()

	postgresConfig, err := config.GetPostgresConfig()
	if err != nil {
		log.Fatalf("не удалось получить конфигурацию Postgres: %s", err.Error())
	}

	db, err := db.NewPostgresClient(ctx, postgresConfig)
	if err != nil {
		log.Fatalf("не удалось подключиться к Postgres: %s", err.Error())
	}

	rabbitmqConfig, err := config.GetRabbitMQConfig()
	if err != nil {
		log.Fatalf("не удалось получить конфигурацию RabbitMQ: %s", err.Error())
	}

	rabbitmqBroker, err := rabbitmq.NewRabbitMQBroker(rabbitmqConfig)
	if err != nil {
		log.Fatalf("не удалось подключиться к RabbitMQ: %s", err.Error())
	}
	// defer rabbitmqBroker.Close()

	app, err := app.NewApp(ctx, db, rabbitmqBroker)
	if err != nil {
		log.Fatalf("ошибка инициализации сервиса: %s", err.Error())
	}

	err = app.RunScadaConsumer(ctx)
	if err != nil {
		log.Fatalf("ошибка запуска сервиса: %v", err)
	}
}
