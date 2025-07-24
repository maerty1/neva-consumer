package zulu

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Данные для получения сырья
type GroupMeasurementsData struct {
	In  string `json:"in"`
	Out string `json:"out"`
}

// Данные из Зулу
type GroupMeasurementsCalculatedData struct {
	In  *float64 `json:"in"`
	Out *float64 `json:"out"`
}

type GroupMeasurement struct {
	ID             int                             `json:"id"`
	Name           string                          `json:"name" example:"Температура"`
	Unit           string                          `json:"unit" example:"атм"`
	CalculatedData GroupMeasurementsCalculatedData `json:"calculated_data"`
	Data           GroupMeasurementsData           `json:"data"`
	Rn             int                             `json:"rn"`
	ZuluCoeff      *float64                        `json:"zulu_coeff"`
	LersCoeff      *float64                        `json:"lers_coeff"`
}

type GetPointDataByCategoryGroup struct {
	Measurements map[string]map[int]*GroupMeasurement `json:"measurements"`
}

func (c *apiClient) GetPointCategoryDataGroup(elemID int, categoryID int, timestamp string) (GetPointDataByCategoryGroup, error) {
	var url string
	if len(timestamp) == 0 {
		url = fmt.Sprintf("%s/zulu/api/v1/points/%v/categories/%v?type=group&n_days=10", c.serviceMapper.GetServiceURL("zulu"), elemID, categoryID)
	} else {
		url = fmt.Sprintf("%s/zulu/api/v1/points/%v/categories/%v?type=group&timestamp=%v", c.serviceMapper.GetServiceURL("zulu"), elemID, categoryID, timestamp)
	}

	req, err := http.NewRequest("GET", url, nil)
	// Хардкодик
	req.Header.Set("X-USER-ID", "1")
	if err != nil {
		return GetPointDataByCategoryGroup{}, fmt.Errorf("ошибка создания запроса: %w", err)
	}
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return GetPointDataByCategoryGroup{}, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GetPointDataByCategoryGroup{}, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var pointData GetPointDataByCategoryGroup
	if err := json.NewDecoder(resp.Body).Decode(&pointData); err != nil {
		return GetPointDataByCategoryGroup{}, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	return pointData, nil
}
