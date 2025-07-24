package zulu

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"
)

type Measurement struct {
	ZuluCoeff      *float64                  `json:"zulu_coeff"`
	LersCoeff      *float64                  `json:"lers_coeff"`
	Name           string                    `json:"name"`
	Unit           string                    `json:"unit"`
	Data           MeasurementData           `json:"data"`
	CalculatedData MeasurementCalculatedData `json:"calculated_data"`
}

type MeasurementData struct {
	In  string `json:"in"`
	Out string `json:"out"`
}

type MeasurementCalculatedData struct {
	In  *float64 `json:"in"`
	Out *float64 `json:"out"`
}

type FullElementData struct {
	Address string                         `json:"address"`
	Title   string                         `json:"title"`
	Packets map[string]map[int]Measurement `json:"packets"`
}

func (c *apiClient) GetFullPoint(elemID int, nDays int) (FullElementData, error) {
	url := fmt.Sprintf("%s/zulu/api/v1/points/%v/full?n_days=%v", c.serviceMapper.GetServiceURL("zulu"), elemID, nDays)

	req, err := http.NewRequest("GET", url, nil)
	// Хардкодик
	req.Header.Set("X-USER-ID", "1")
	if err != nil {
		return FullElementData{}, fmt.Errorf("ошибка создания запроса: %w", err)
	}
	client := &http.Client{}
	client.Timeout = 600 * time.Second

	resp, err := client.Do(req)
	if err != nil {
		return FullElementData{}, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return FullElementData{}, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var fullPoint FullElementData
	if err := json.NewDecoder(resp.Body).Decode(&fullPoint); err != nil {
		return FullElementData{}, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	// Применяем корректировки
	adjustCalculatedData(&fullPoint, 3)

	return fullPoint, nil
}
func adjustCalculatedData(data *FullElementData, targetGroupID int) {
	for timestamp, groups := range data.Packets {
		for groupID, measurement := range groups {
			if groupID == targetGroupID {
				if measurement.CalculatedData.In != nil {
					*measurement.CalculatedData.In = roundTo(*measurement.CalculatedData.In*24, 2)
				}
				if measurement.CalculatedData.Out != nil {
					*measurement.CalculatedData.Out = roundTo(*measurement.CalculatedData.Out*24, 2)
				}

				groups[groupID] = measurement
			}
		}

		data.Packets[timestamp] = groups
	}
}
func roundTo(value float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Round(value*pow) / pow
}
