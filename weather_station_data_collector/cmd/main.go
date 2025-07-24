package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"weather_station_data_collector/internal/app"
	"weather_station_data_collector/internal/config"
	"weather_station_data_collector/internal/db"
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

	app.RunInterrogator(ctx)
}
