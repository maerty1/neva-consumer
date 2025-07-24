package lers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type MeasurePoint struct {
	ID         int
	DeviceID   int    `json:"deviceId"`
	FullTitle  string `json:"fullTitle"`
	Title      string `json:"title"`
	Address    string `json:"address"`
	SystemType string `json:"systemType"`
}

type MeasurePointsResponse struct {
	MeasurePoints []MeasurePoint `json:"measurePoints"`
}

// GetMeasurePoints возвращает список точек измерения
func (c *apiClient) GetMeasurePoints(token string, serverHost string) ([]MeasurePoint, error) {
	url := fmt.Sprintf("%s/api/v1/Core/MeasurePoints", serverHost)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("не удалось получить точки измерения: %s", resp.Status)
	}

	var response MeasurePointsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.MeasurePoints, nil
}
