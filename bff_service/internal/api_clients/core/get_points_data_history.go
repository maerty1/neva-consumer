package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
)

type GetPointsDataHistoryResponse map[string]map[string]map[string]GetPointsDataHistoryMeasurementsResponse

type GetPointsDataHistoryMeasurementsRequest struct {
	I string `json:"i"`
	O string `json:"o"`
}

type GetPointsDataHistoryMeasurementsResponse struct {
	I *float64 `json:"i"`
	O *float64 `json:"o"`
}

type GetPointsDataHistoryRequest struct {
	ElemID       int                                             `json:"elem_id"`
	Measurements map[int]GetPointsDataHistoryMeasurementsRequest `json:"measurements"`
}

func (a *apiClient) GetPointsDataHistory(reqData []GetPointsDataHistoryRequest, nDays int, timestamp string) (GetPointsDataHistoryResponse, error) {
	url := fmt.Sprintf("%s/core/api/v1/points/data?n_days=%v", a.serviceMapper.GetServiceURL("core"), nDays)

	jsonPayload, err := json.Marshal(reqData)
	if err != nil {
		return GetPointsDataHistoryResponse{}, fmt.Errorf("ошибка маршалинга payload: %w", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return GetPointsDataHistoryResponse{}, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("X-USER-ID", "1")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GetPointsDataHistoryResponse{}, fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GetPointsDataHistoryResponse{}, fmt.Errorf("не удалось авторизироваться: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	var response GetPointsDataHistoryResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return GetPointsDataHistoryResponse{}, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}
	// convertAtmospheresToMeters(&response, []string{"2"})

	return response, nil
}
func convertAtmospheresToMeters(data *GetPointsDataHistoryResponse, targetGroupIDs []string) {
	const conversionFactor = 10.33227

	for elemID, timestamps := range *data {
		for timestamp, groups := range timestamps {
			for groupID, measurements := range groups {
				if contains(targetGroupIDs, groupID) {
					if measurements.I != nil {
						*measurements.I = roundTo(*measurements.I*conversionFactor, 2)
					}
					if measurements.O != nil {
						*measurements.O = roundTo(*measurements.O*conversionFactor, 2)
					}
					groups[groupID] = measurements
				}
			}
			timestamps[timestamp] = groups
		}
		(*data)[elemID] = timestamps
	}
}

// contains проверяет, содержит ли срез strSlice строку str
func contains(strSlice []string, str string) bool {
	for _, v := range strSlice {
		if v == str {
			return true
		}
	}
	return false
}

// roundTo округляет число до заданного количества десятичных знаков
func roundTo(value float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Round(value*pow) / pow
}
