package zulu

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MeasurementKeyvalue struct {
	Name   string      `json:"name" example:"Температура"`
	Unit   string      `json:"unit" example:"атм"`
	Source string      `json:"source" example:"zulu/scada"`
	Value  interface{} `json:"value"`
	Rn     int         `json:"rn"`
}
type GetPointDataByCategoryKeyvalue struct {
	Measurements []MeasurementKeyvalue `json:"measurements"`
}

func (c *apiClient) GetPointCategoryDataKeyvalue(elemID int, categoryID int) (GetPointDataByCategoryKeyvalue, error) {
	url := fmt.Sprintf("%s/zulu/api/v1/points/%v/categories/%v?type=keyvalue&n_days=10", c.serviceMapper.GetServiceURL("zulu"), elemID, categoryID)

	req, err := http.NewRequest("GET", url, nil)
	// Хардкодик
	req.Header.Set("X-USER-ID", "1")
	if err != nil {
		return GetPointDataByCategoryKeyvalue{}, fmt.Errorf("ошибка создания запроса: %w", err)
	}
	client := &http.Client{}
	client.Timeout = 600 * time.Second

	resp, err := client.Do(req)
	if err != nil {
		return GetPointDataByCategoryKeyvalue{}, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GetPointDataByCategoryKeyvalue{}, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var pointData GetPointDataByCategoryKeyvalue
	if err := json.NewDecoder(resp.Body).Decode(&pointData); err != nil {
		return GetPointDataByCategoryKeyvalue{}, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	return pointData, nil
}
