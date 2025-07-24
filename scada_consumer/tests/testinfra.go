package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"scada_consumer/internal/app"
	"scada_consumer/internal/config"
	"scada_consumer/internal/db"
	"scada_consumer/internal/message_broker/rabbitmq"
	"sync"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	rabbitTestContainer "github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	AppInstance *app.App
	CleanupFunc func()
	once        sync.Once
	setupErr    error
)

func setup(ctx context.Context) {
	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	testsBaseDir := os.Getenv("TESTS_BASE_DIR")

	if testsBaseDir == "" {
		setupErr = fmt.Errorf("не установлена переменная окружения TESTS_BASE_DIR")
		return
	}

	postgresContainer, err := postgres.Run(ctx,
		"docker.io/postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.WithInitScripts(filepath.Join(testsBaseDir, "/data/init-db.sh")),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		setupErr = fmt.Errorf("не удалось запустить контейнер Postgres: %w", err)
		return
	}

	rabbitContainer, err := rabbitTestContainer.Run(ctx,
		"rabbitmq:3.12.11-management-alpine",
		// rabbitTestContainer.WithAdminUsername("admin"),
		// rabbitTestContainer.WithAdminPassword("password"),
	)
	if err != nil {
		setupErr = fmt.Errorf("не удалось запустить контейнер RabbitMQ: %w", err)
		return
	}

	rabbitUrl, err := rabbitContainer.AmqpURL(ctx)
	if err != nil {
		log.Fatal(err)
	}

	os.Setenv("RABBITMQ_URL", rabbitUrl)
	// os.Setenv("RABBITMQ_QUEUE", "test-queue")

	pgHost, err := postgresContainer.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	pgPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatal(err)
	}

	pgDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, pgHost, pgPort.Port(), dbName)
	pgConfig := config.NewPGConfig(pgDSN, 1, 30)
	rabbitmqConfig, err := config.GetRabbitMQConfig()
	if err != nil {
		setupErr = fmt.Errorf("ошибка получения конфигурации RabbitMQ: %w", err)
		return
	}

	postgresClient, err := db.NewPostgresClient(ctx, pgConfig)
	if err != nil {
		setupErr = fmt.Errorf("ошибка инициализации клиента Postgres: %w", err)
		return
	}

	rabbitmqBroker, err := rabbitmq.NewRabbitMQBroker(rabbitmqConfig)
	if err != nil {
		setupErr = fmt.Errorf("ошибка получения брокера Rabbitmq: %w", err)
		return
	}

	AppInstance, err = app.NewApp(ctx, postgresClient, rabbitmqBroker)
	if err != nil {
		setupErr = fmt.Errorf("ошибка инициализации приложения: %w", err)
		return
	}

	CleanupFunc = func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("не удалось завершить работу контейнера Postgres: %s", err)
		}
		if err := rabbitContainer.Terminate(ctx); err != nil {
			log.Fatalf("не удалось завершить контейнер RabbitMQ: %s", err)
		}
	}
}

func GetApp() (*app.App, func(), error) {
	return AppInstance, CleanupFunc, setupErr
}

func Init(ctx context.Context) {
	once.Do(func() {
		setup(ctx)
	})
}
