package weather_station_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"weather_station_data_collector/internal/api_client/weather_station"
	"weather_station_data_collector/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestGetHistoryMinute_Success(t *testing.T) {
	expectedResponse := []models.WeatherDataRaw{
		{
			Dateutc:        "2024-11-25 12:00:00",
			Tempinf:        "72.5",
			Humidityin:     "45",
			Baromrelin:     "29.92",
			Baromabsin:     "29.85",
			Tempf:          "68.0",
			Humidity:       "50",
			Winddir:        "180",
			Windspeedmph:   "5.0",
			Windgustmph:    "7.5",
			Maxdailygust:   "15.0",
			Solarradiation: "120.5",
			Uv:             "3",
			Rainratein:     "0.01",
			Eventrainin:    "0.1",
			Hourlyrainin:   "0.02",
			Dailyrainin:    "0.25",
			Weeklyrainin:   "0.5",
			Monthlyrainin:  "1.2",
			Yearlyrainin:   "10.0",
			Totalrainin:    "50.0",
			Wh65Batt:       "1",
			Freq:           "868M",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/weather/pull/history", r.URL.Path[:21]) // Проверяем путь
		assert.Contains(t, r.URL.RawQuery, "from=2024-01-01")     // Проверяем параметр `from`
		assert.Contains(t, r.URL.RawQuery, "to=2024-01-02")       // Проверяем параметр `to`

		w.WriteHeader(http.StatusOK)
		respBytes, _ := json.Marshal(expectedResponse)
		w.Write(respBytes)
	}))
	defer ts.Close()

	os.Setenv("WEATHER_STATION_BASE_URL", ts.Listener.Addr().String())
	os.Setenv("NETWORK_PROTOCOL", "http")

	client := weather_station.NewApiClient()

	ctx := context.Background()
	from := "2024-01-01"
	to := "2024-01-02"
	result, err := client.GetHistoryMinute(ctx, from, to)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
}

func TestGetHistoryMinute_RequestError(t *testing.T) {
	os.Setenv("WEATHER_STATION_BASE_URL", "invalid-url")
	os.Setenv("NETWORK_PROTOCOL", "http")

	client := weather_station.NewApiClient()

	ctx := context.Background()
	from := "2024-01-01"
	to := "2024-01-02"
	result, err := client.GetHistoryMinute(ctx, from, to)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetHistoryMinute_UnsuccessfulStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	os.Setenv("WEATHER_STATION_BASE_URL", ts.Listener.Addr().String())
	os.Setenv("NETWORK_PROTOCOL", "http")

	client := weather_station.NewApiClient()

	ctx := context.Background()
	from := "2024-01-01"
	to := "2024-01-02"
	result, err := client.GetHistoryMinute(ctx, from, to)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetHistoryMinute_InvalidJSONResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid-json"))
	}))
	defer ts.Close()

	os.Setenv("WEATHER_STATION_BASE_URL", ts.Listener.Addr().String())
	os.Setenv("NETWORK_PROTOCOL", "http")

	client := weather_station.NewApiClient()

	ctx := context.Background()
	from := "2024-01-01"
	to := "2024-01-02"
	result, err := client.GetHistoryMinute(ctx, from, to)

	assert.Error(t, err)
	assert.Nil(t, result)
}
