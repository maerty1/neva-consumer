package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"zulu_updater/internal/app"
	"zulu_updater/internal/config"
	"zulu_updater/internal/db"
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

	rootPostgresConfig, err := config.GetPostgresConfig("root")
	if err != nil {
		log.Fatalf("не удалось получить конфигурацию Postgres: %s", err.Error())
	}

	zuluPostgresConfig, err := config.GetPostgresConfig("zulu")
	if err != nil {
		log.Fatalf("не удалось получить конфигурацию Postgres: %s", err.Error())
	}

	rootDb, err := db.NewPostgresClient(ctx, rootPostgresConfig)
	if err != nil {
		log.Fatalf("не удалось подключиться к Postgres: %s", err.Error())
	}

	zuluDb, err := db.NewPostgresClient(ctx, zuluPostgresConfig)

	app, err := app.NewApp(ctx, rootDb, zuluDb)
	if err != nil {
		log.Fatalf("ошибка инициализации сервиса: %s", err.Error())
	}

	app.RunZuluUpdater(ctx)
}
