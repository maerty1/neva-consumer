package weather_station

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"weather_station_data_collector/internal/models"
)

func (a *apiClient) GetLastData(ctx context.Context) models.WeatherStationResponseChannel {
	url := fmt.Sprintf("%s/weather/pull/last", a.baseUrl)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res := models.WeatherStationResponseChannel{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		res.Err = err
		return res
	}

	resp, err := a.client.Do(req.WithContext(ctx))
	if err != nil {
		res.Err = err
		return res
	}
	defer resp.Body.Close()

	res.Code = resp.StatusCode

	if resp.StatusCode != http.StatusOK {
		return res
	}

	var response models.WeatherDataRaw
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		res.Err = err
		return res
	}
	res.Body = response

	return res
}
