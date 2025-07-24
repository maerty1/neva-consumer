package forecast_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"weather_station_data_collector/internal/api_client/forecast"
	"weather_station_data_collector/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestGetForecast_Success(t *testing.T) {
	expectedResponse := models.WeatherResponse{
		Cod:     "200",
		Message: 0,
		Cnt:     1,
		List: []models.WeatherItem{
			{
				Dt: 1633046400,
				Main: models.Main{
					Temp:      15.5,
					FeelsLike: 15.0,
					TempMin:   14.0,
					TempMax:   16.0,
					Pressure:  1012,
					Humidity:  78,
				},
				Weather: []models.WeatherDetail{
					{
						ID:          800,
						Main:        "Clear",
						Description: "clear sky",
						Icon:        "01d",
					},
				},
				Clouds: models.Clouds{
					All: 0,
				},
				Wind: models.Wind{
					Speed: 3.5,
					Deg:   120,
					Gust:  4.0,
				},
				Rain: nil,
				Sys: models.Sys{
					Pod: "d",
				},
				DtTxt: "2024-11-29 12:00:00",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/data/2.5/forecast")
		assert.Contains(t, r.URL.RawQuery, "lat=59.11726")
		assert.Contains(t, r.URL.RawQuery, "lon=28.086979")
		assert.Contains(t, r.URL.RawQuery, "units=metric")

		w.WriteHeader(http.StatusOK)
		respBytes, _ := json.Marshal(expectedResponse)
		w.Write(respBytes)
	}))
	defer ts.Close()

	os.Setenv("FORECAST_BASE_URL", "http://"+ts.Listener.Addr().String())
	os.Setenv("FORECAST_API_TOKEN", "mockToken")

	client := forecast.NewApiClient()

	ctx := context.Background()
	result, err := client.GetForecast(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, *result)
}

func TestGetForecast_RequestError(t *testing.T) {
	os.Setenv("FORECAST_BASE_URL", "invalid-url")
	os.Setenv("FORECAST_API_TOKEN", "mockToken")

	client := forecast.NewApiClient()

	ctx := context.Background()
	result, err := client.GetForecast(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetForecast_UnsuccessfulStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	os.Setenv("FORECAST_BASE_URL", ts.Listener.Addr().String())
	os.Setenv("FORECAST_API_TOKEN", "mockToken")

	client := forecast.NewApiClient()

	ctx := context.Background()
	result, err := client.GetForecast(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "код ответа:500")
}

func TestGetForecast_InvalidJSONResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid-json"))
	}))
	defer ts.Close()

	os.Setenv("FORECAST_BASE_URL", "http://"+ts.Listener.Addr().String())
	os.Setenv("FORECAST_API_TOKEN", "mockToken")

	client := forecast.NewApiClient()

	ctx := context.Background()
	result, err := client.GetForecast(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid character")
}
