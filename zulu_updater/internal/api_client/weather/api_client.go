package weather

import (
	"net/http"
	"os"
	"zulu_updater/internal/models"
)

type ApiClient interface {
	GetWeatherData() (*models.WeatherData, error)
}

var _ ApiClient = (*apiClient)(nil)

type apiClient struct {
	baseUrl string
	client  *http.Client
}

func NewApiClient() ApiClient {
	var res apiClient
	res.baseUrl = os.Getenv("WEATHER_BASE_URL")
	res.client = &http.Client{}
	return &res
}
