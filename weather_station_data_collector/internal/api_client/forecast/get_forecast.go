package forecast

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"weather_station_data_collector/internal/models"
)

func (a *apiClient) GetForecast(ctx context.Context) (*models.WeatherResponse, error) {
	url := fmt.Sprintf("%s/data/2.5/forecast?lat=%s&lon=%s&appid=%s&units=metric",
		a.baseUrl, a.lat, a.lon, a.apiToken)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("код ответа:%d", resp.StatusCode))
	}

	var response models.WeatherResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
