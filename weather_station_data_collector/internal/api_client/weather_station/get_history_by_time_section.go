package weather_station

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"weather_station_data_collector/internal/models"

	"golang.org/x/net/context"
)

func (a *apiClient) GetHistoryMinute(ctx context.Context, from string, to string) ([]models.WeatherDataRaw, error) {
	url := fmt.Sprintf("%s/weather/pull/history?from=%s&to=%s", a.baseUrl, from, to)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := a.client.Do(req.WithContext(ctx))
	if err != nil {
		log.Println("Error performing request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error response for %s: %d\n", url, resp.StatusCode)
	}

	var response []models.WeatherDataRaw
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}
