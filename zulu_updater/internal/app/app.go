package app

import (
	"context"
	"log"
	"zulu_updater/internal/api_client/weather"
	zuluApi "zulu_updater/internal/api_client/zulu"
	"zulu_updater/internal/db"
	zuluUpdater "zulu_updater/internal/services/zulu_updater"
	measurePointsDataDay "zulu_updater/repositories/measure_points_data_day"
	zuluRepo "zulu_updater/repositories/zulu"
)

type App struct {
	zuluUpdater zuluUpdater.Service
}

func NewApp(ctx context.Context, rootDb db.PostgresClient, zuluDb db.PostgresClient) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx, rootDb, zuluDb)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) RunZuluUpdater(ctx context.Context) {
	a.runZuluUpdater(ctx)
}

func (a *App) initDeps(ctx context.Context, rootDb db.PostgresClient, zuluDb db.PostgresClient) error {
	a.zuluUpdater = zuluUpdater.NewService(
		measurePointsDataDay.NewRepository(rootDb),
		zuluApi.NewApiClient(),
		weather.NewApiClient(),
		zuluRepo.NewRepository(zuluDb),
	)
	return nil
}

func (a *App) runZuluUpdater(ctx context.Context) {
	log.Printf("Запуск обновления")
	a.zuluUpdater.RunZuluUpdater(ctx)
}
