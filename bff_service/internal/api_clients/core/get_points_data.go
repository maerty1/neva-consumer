package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GetPointsDataRequestMeasurement struct {
	I string `json:"i" example:"T_in"`
	O string `json:"o" example:"T_out"`
}

type GetPointsDataRequest struct {
	ElemID       int                                        `json:"elem_id"`
	Measurements map[string]GetPointsDataRequestMeasurement `json:"measurements"`
}

type GetPointsDataResponseMeasurement struct {
	I *float64 `json:"i" example:"53.24"`
	O *float64 `json:"o" example:"45.12"`
}

type GetPointsDataResponse struct {
	ElemID       int                                         `json:"elem_id"`
	IsDataCopied bool                                        `json:"iscopied"`
	Measurements map[string]GetPointsDataResponseMeasurement `json:"measurements"`
}

func (c *apiClient) GetPointsData(reqData []GetPointsDataRequest, timestamp string) ([]GetPointsDataResponse, error) {
	var url string

	if len(timestamp) == 0 {
		url = fmt.Sprintf("%s/core/api/v1/points/last_data", c.serviceMapper.GetServiceURL("core"))
	} else {
		url = fmt.Sprintf("%s/core/api/v1/points/last_data?timestamp=%v", c.serviceMapper.GetServiceURL("core"), timestamp)
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации данных запроса: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("X-USER-ID", "1")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("неожиданный статус-код: %d. Тело ответа: %s", resp.StatusCode, string(bodyBytes))
	}

	var pointsData []GetPointsDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&pointsData); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	return pointsData, nil
}
