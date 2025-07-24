package weather_data_raw

import (
	"context"
	"time"
	"weather_station_data_collector/internal/db"
	"weather_station_data_collector/internal/models"
)

type Repository interface {
	InsertRawWeatherData(ctx context.Context, rawData models.WeatherDataRaw) error
	SelectTimeWithNullData(ctx context.Context) ([]time.Time, error)
	UpdateRawWeatherData(ctx context.Context, rawData models.WeatherDataRaw, oldTime time.Time) error
	InsertWeatherDataBatch(ctx context.Context, data []models.WeatherDataRaw) error

	SelectWeatherDataByDay(ctx context.Context, day time.Time) ([]models.WeatherDataRaw, float64, float64, error)
}

type repository struct {
	db db.PostgresClient
}

func NewRepository(db db.PostgresClient) *repository {
	return &repository{
		db: db,
	}
}
