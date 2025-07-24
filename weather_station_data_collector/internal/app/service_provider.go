package app

import (
	weatherStation "weather_station_data_collector/internal/api_client/weather_station"
	"weather_station_data_collector/internal/db"
	weatherDataRaw "weather_station_data_collector/internal/repositories/weather_data_raw"
)

type serviceProvider struct {
	postgresDB               db.PostgresClient
	weatherDataRawRepository weatherDataRaw.Repository
	weatherStationApiClient  weatherStation.ApiClient
}

func newServiceProvider(db db.PostgresClient) *serviceProvider {
	return &serviceProvider{
		postgresDB: db,
	}
}

func (s *serviceProvider) PostgresDB() db.PostgresClient {
	return s.postgresDB
}

func (s *serviceProvider) WeatherDataRawRepository() weatherDataRaw.Repository {
	if s.weatherDataRawRepository == nil {
		s.weatherDataRawRepository = weatherDataRaw.NewRepository(s.PostgresDB())
	}
	return s.weatherDataRawRepository
}

func (s *serviceProvider) WeatherStationApiClient() weatherStation.ApiClient {
	if s.weatherStationApiClient == nil {
		s.weatherStationApiClient = weatherStation.NewApiClient()
	}
	return s.weatherStationApiClient
}
