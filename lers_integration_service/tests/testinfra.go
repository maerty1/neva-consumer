package tests

import (
	"context"
	"fmt"
	"lers_integration_service/internal/app"
	"lers_integration_service/internal/config"
	"lers_integration_service/internal/db"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/testcontainers/testcontainers-go"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
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
		"postgres:16-alpine",
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
	if err != nil {
		setupErr = fmt.Errorf("ошибка получения конфигурации Postgres: %w", err)
		return
	}

	postgresClient, err := db.NewPostgresClient(ctx, pgConfig)
	if err != nil {
		setupErr = fmt.Errorf("ошибка инициализации клиента Postgres: %w", err)
		return
	}

	AppInstance, err = app.NewApp(ctx, postgresClient)
	if err != nil {
		setupErr = fmt.Errorf("ошибка инициализации приложения: %w", err)
		return
	}

	CleanupFunc = func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("не удалось завершить работу контейнера Postgres: %s", err)
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

func CleanDb(ctx context.Context, app *app.App) error {
	_, err := app.S().PostgresDB().DB().Exec(ctx, "TRUNCATE accounts, measure_points, measure_points_data, accounts_sync_log, measure_points_poll_retry")
	if err != nil {
		return fmt.Errorf("ошибка очистки БД: %w", err)
	}

	return nil
}
