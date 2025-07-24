package weather_station

import (
	"context"
	"net/http"
	"os"
	"weather_station_data_collector/internal/models"
)

type ApiClient interface {
	GetLastData(ctx context.Context) models.WeatherStationResponseChannel
	GetHistoryMinute(ctx context.Context, from string, to string) ([]models.WeatherDataRaw, error)
}

var _ ApiClient = (*apiClient)(nil)

type apiClient struct {
	baseUrl string
	client  *http.Client
}

func NewApiClient() ApiClient {
	return &apiClient{
		baseUrl: os.Getenv("WEATHER_STATION_BASE_URL"),
		client:  &http.Client{},
	}
}
