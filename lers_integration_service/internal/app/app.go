package app

import (
	"context"
	"lers_integration_service/internal/db"
	"log"
)

type App struct {
	serviceProvider *serviceProvider
}

func NewApp(ctx context.Context, db db.PostgresClient) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx, db)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) RunSynchronizer(ctx context.Context) error {
	return a.runSynchronizer(ctx)
}

func (a *App) RunPoller(ctx context.Context) error {
	return a.runPoller(ctx)
}

func (a *App) RunRetryer(ctx context.Context) error {
	return a.runRetryer(ctx)
}

func (a *App) initDeps(ctx context.Context, db db.PostgresClient) error {
	a.initServiceProvider(ctx, db)

	return nil
}

func (a *App) initServiceProvider(_ context.Context, db db.PostgresClient) error {
	a.serviceProvider = newServiceProvider(db)
	return nil
}

func (a *App) runSynchronizer(ctx context.Context) error {
	log.Printf("Запуск синхронизации")
	err := a.serviceProvider.SynchronizerService().RunSynchronizer(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) runPoller(ctx context.Context) error {
	log.Printf("Запуск опроса")
	err := a.serviceProvider.PollerService().RunPoller(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) runRetryer(ctx context.Context) error {
	log.Printf("Запуск retryer")
	err := a.serviceProvider.RetryerService().RunRetryer(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) S() *serviceProvider {
	return a.serviceProvider
}
