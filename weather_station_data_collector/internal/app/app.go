package app

import (
	"context"
	"log"
	"weather_station_data_collector/internal/api_client/forecast"
	"weather_station_data_collector/internal/db"
	weatherData "weather_station_data_collector/internal/repositories/weather_data"
	weatherDataRaw "weather_station_data_collector/internal/repositories/weather_data_raw"
	dayAvg "weather_station_data_collector/internal/services/day_avg"
	"weather_station_data_collector/internal/services/interrogator"
)

type App struct {
	serviceProvider *serviceProvider
	interrogator    interrogator.Service
	avgCalculator   dayAvg.Service
}

func NewApp(ctx context.Context, db db.PostgresClient) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx, db)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) RunInterrogator(ctx context.Context) {
	a.runInterrogator(ctx)
}

func (a *App) RunHandleInterrogator(ctx context.Context, timeFrom, timeTo string) {
	a.runHandleInterrogator(ctx, timeFrom, timeTo)
}

func (a *App) RunAvgCalculator(ctx context.Context) {
	a.runAvgCalculator(ctx)
}

func (a *App) initDeps(ctx context.Context, db db.PostgresClient) error {
	a.initServiceProvider(ctx, db)

	return nil
}

func (a *App) initServiceProvider(_ context.Context, db db.PostgresClient) error {
	a.serviceProvider = newServiceProvider(db)
	a.interrogator = interrogator.NewService(
		weatherDataRaw.NewRepository(db),
		a.serviceProvider.WeatherStationApiClient(),
	)
	a.avgCalculator = dayAvg.NewService(
		weatherDataRaw.NewRepository(db),
		weatherData.NewRepository(db),
		forecast.NewApiClient(),
	)
	return nil
}

func (a *App) runInterrogator(ctx context.Context) {
	log.Printf("Запуск опроса")
	a.interrogator.RunInterrogator(ctx)
}

func (a *App) runHandleInterrogator(ctx context.Context, timeFrom, timeTo string) {
	log.Printf("Запуск опроса")
	a.interrogator.RunHandleInterrogator(ctx, timeFrom, timeTo)
}

func (a *App) runAvgCalculator(ctx context.Context) {
	log.Printf("Запуск расчётов")
	a.avgCalculator.CalculateCurrentDayAvg(ctx)
	log.Printf("Заполнение прогноза")
	a.avgCalculator.CalculateForecastAvg(ctx)
}
