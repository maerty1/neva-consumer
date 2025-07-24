package zulu

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type MeasurementGroup struct {
	Coeff *float64 `json:"coeff"`
	In    string   `json:"i"`
	Out   string   `json:"o"`
}

type Point struct {
	ElemID            int                      `json:"elem_id"`
	Title             string                   `json:"title" example:"Котельная 22"`
	Address           string                   `json:"address" example:"Улица Пушкина 12"`
	MeasurementGroups map[int]MeasurementGroup `json:"measurement_groups"`
	Coordinates       []float64                `json:"coordinates" example:"[45.61888, 75.35849]" description:"[lat, lon]"`
	HasAccident       bool                     `json:"has_accident"`
	IsCopied          bool                     `json:"iscopied"`
	Type              int                      `json:"type"`
}

func (c *apiClient) GetPoints(zwsTypeIds []int) ([]Point, error) {
	zwsTypeIdsStr := strings.Trim(strings.Replace(fmt.Sprint(zwsTypeIds), " ", ",", -1), "[]")
	url := fmt.Sprintf("%s/zulu/api/v1/points?zws_type_id=%s", c.serviceMapper.GetServiceURL("zulu"), zwsTypeIdsStr)

	req, err := http.NewRequest("GET", url, nil)
	// Хардкодик
	req.Header.Set("X-USER-ID", "1")
	if err != nil {
		return []Point{}, fmt.Errorf("ошибка создания запроса: %w", err)
	}
	client := &http.Client{}
	client.Timeout = 600 * time.Second

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var points []Point
	if err := json.NewDecoder(resp.Body).Decode(&points); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	return points, nil
}
