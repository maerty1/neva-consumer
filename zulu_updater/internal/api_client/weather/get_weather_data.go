package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"zulu_updater/internal/models"
)

func (a apiClient) GetWeatherData() (*models.WeatherData, error) {
	url := fmt.Sprintf("%s/core/api/v1/weather/with_forecast",
		a.baseUrl)
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

	var response models.WeatherData
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
