package zulu

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CalculatedData struct {
	In  *float64 `json:"in"`
	Out *float64 `json:"out"`
}

type Data struct {
	In  *string `json:"in"`
	Out *string `json:"out"`
}

type GetPointsDataCategoryMeasurement struct {
	Name           string         `json:"name"`
	Unit           string         `json:"unit"`
	CalculatedData CalculatedData `json:"calculated_data"`
	Data           Data           `json:"data"`
	Rn             int            `json:"rn"`
	ZuluCoeff      *float64       `json:"zulu_coeff"`
	LersCoeff      *float64       `json:"lers_coeff"`
}

type CategoryMeasurements struct {
	Measurements       map[string]GetPointsDataCategoryMeasurement `json:"measurements"`
	IsCalculatedCopied bool                                        `json:"iscopied"`
}

// Тип для всего ответа, где ключ — ID категории
type GetPointsDataCategoryResponse map[string]CategoryMeasurements

func (c *apiClient) GetPointsDataCategoryByZwsType(zwsTypeIds []int, categoryID int, timestamp string) (GetPointsDataCategoryResponse, error) {
	zwsTypeIdsStr := strings.Trim(strings.Replace(fmt.Sprint(zwsTypeIds), " ", ",", -1), "[]")
	url := fmt.Sprintf("%s/zulu/api/v2/points/categories/%v?zws_type_id=%v&timestamp=%v", c.serviceMapper.GetServiceURL("zulu"), categoryID, zwsTypeIdsStr, timestamp)

	req, err := http.NewRequest("GET", url, nil)
	// Хардкодик
	req.Header.Set("X-USER-ID", "1")
	if err != nil {
		return GetPointsDataCategoryResponse{}, fmt.Errorf("ошибка создания запроса: %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var points GetPointsDataCategoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&points); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	return points, nil
}
