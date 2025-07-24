package lers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ConsumptionData struct {
	ResourceKind string  `json:"resourceKind"`
	DateTime     string  `json:"dateTime"`
	Values       []Value `json:"values"`
}

type Value struct {
	DataParameter  string  `json:"dataParameter"`
	Value          float64 `json:"value"`
	IsBad          bool    `json:"isBad"`
	IsCalc         bool    `json:"isCalc"`
	IsInterpolated bool    `json:"isInterpolated"`
	IsReset        bool    `json:"isReset"`
}

type ConsumptionResponse struct {
	PressureType    string            `json:"pressureType"`
	HourConsumption []ConsumptionData `json:"hourConsumption"`
	DayConsumption  []ConsumptionData `json:"dayConsumption"`
}

// GetConsumptionData это основной метод для получения исторических данных от точки измерения
func (c *apiClient) GetConsumptionData(accountID int, token string, serverHost string, measurePointID int, startDate, endDate string) (*ConsumptionResponse, error) {
	url := fmt.Sprintf("%s/api/v1/Data/MeasurePoints/%d/Consumption/%s/%s?dataTypes=Day&electricDataKind=Raw&includeCalculated=true&includeAbsentRecords=true&withSummary=true&considerReportingDate=true&units=ConfiguredUnits", serverHost, measurePointID, startDate, endDate)
	// url := fmt.Sprintf("%s/api/v1/Data/MeasurePoints/%d/Consumption/%s/%s?dataTypes=Hour&dataTypes=Day&electricDataKind=Raw&includeCalculated=true&includeAbsentRecords=true&withSummary=true&considerReportingDate=true&units=ConfiguredUnits", serverHost, measurePointID, startDate, endDate)
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
		return nil, fmt.Errorf("не удалось получить данные о потреблении: %s", resp.Status)
	}

	var response ConsumptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
