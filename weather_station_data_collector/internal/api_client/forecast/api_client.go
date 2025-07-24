package forecast

import (
	"context"
	"net/http"
	"os"
	"weather_station_data_collector/internal/models"
)

type ApiClient interface {
	GetForecast(ctx context.Context) (*models.WeatherResponse, error)
}

var _ ApiClient = (*apiClient)(nil)

type apiClient struct {
	baseUrl  string
	apiToken string
	client   *http.Client
	lat      string
	lon      string
}

func NewApiClient() ApiClient {
	var res apiClient
	if os.Getenv("LAT_CORD") == "" {
		res.lat = "59.11726"
		res.lon = "28.086979"
	} else {
		res.lat = os.Getenv("LAT_CORD")
		res.lon = os.Getenv("LON_CORD")
	}
	res.baseUrl = os.Getenv("FORECAST_BASE_URL")
	res.apiToken = os.Getenv("FORECAST_API_TOKEN")
	res.client = &http.Client{}
	return &res
}
