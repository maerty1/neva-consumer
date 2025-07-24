package day_avg

import (
	"context"
	"weather_station_data_collector/internal/api_client/forecast"
	weatherData "weather_station_data_collector/internal/repositories/weather_data"
	weatherDataRaw "weather_station_data_collector/internal/repositories/weather_data_raw"
)

type Service interface {
	CalculateCurrentDayAvg(ctx context.Context)
	CalculateForecastAvg(ctx context.Context)
}

var _ Service = (*service)(nil)

type service struct {
	RawWeatherDataRepository weatherDataRaw.Repository
	WeatherDataRepository    weatherData.Repository
	ForecastApiClient        forecast.ApiClient
}

func NewService(repositoryRaw weatherDataRaw.Repository, repository weatherData.Repository, apiClient forecast.ApiClient) Service {
	return &service{
		RawWeatherDataRepository: repositoryRaw,
		WeatherDataRepository:    repository,
		ForecastApiClient:        apiClient,
	}
}
