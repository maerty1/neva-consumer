package weather_data

import (
	"context"
	"weather_station_data_collector/internal/db"
	"weather_station_data_collector/internal/models"
)

type Repository interface {
	InsertWeatherData(ctx context.Context, data models.WeatherData) error
	InsertWeatherDataButch(ctx context.Context, data []models.WeatherData) error
}

type repository struct {
	db db.PostgresClient
}

func NewRepository(db db.PostgresClient) *repository {
	return &repository{
		db: db,
	}
}
