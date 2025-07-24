package main

import (
	"context"
	"lers_integration_service/internal/app"
	"lers_integration_service/internal/config"
	"lers_integration_service/internal/db"
	"log"
	"os"
	"os/signal"

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

	app, err := app.NewApp(ctx, db)
	if err != nil {
		log.Fatalf("ошибка инициализации сервиса: %s", err.Error())
	}

	err = app.RunPoller(ctx)
	if err != nil {
		log.Fatalf("ошибка запуска Poller: %v", err)
	}
	err = app.RunRetryer(ctx)
	if err != nil {
		log.Fatalf("ошибка запуска Retryer: %v", err)
	}
	err = app.RunSynchronizer(ctx)
	if err != nil {
		log.Fatalf("ошибка запуска Synchronizer: %v", err)

	}
}
