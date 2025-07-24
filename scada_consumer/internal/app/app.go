package app

import (
	"context"
	"log"
	"scada_consumer/internal/db"
	"scada_consumer/internal/message_broker"
)

type App struct {
	serviceProvider *serviceProvider
}

func NewApp(ctx context.Context, db db.PostgresClient, rabbitBroker message_broker.MessageBroker) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx, db, rabbitBroker)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) RunScadaConsumer(ctx context.Context) error {
	return a.runScadaConsumer(ctx)
}

func (a *App) initDeps(ctx context.Context, db db.PostgresClient, rabbitBroker message_broker.MessageBroker) error {
	a.initServiceProvider(ctx, db, rabbitBroker)

	return nil
}

func (a *App) initServiceProvider(_ context.Context, db db.PostgresClient, rabbitBroker message_broker.MessageBroker) error {
	a.serviceProvider = newServiceProvider(db, rabbitBroker)
	return nil
}

func (a *App) runScadaConsumer(ctx context.Context) error {
	log.Printf("Scada consumer запущен...")
	err := a.serviceProvider.MessageService().RunScadaConsumer(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) S() *serviceProvider {
	return a.serviceProvider
}
