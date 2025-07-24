package zulu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type req struct {
	IDs []int `json:"ids"`
}

func (c *apiClient) GetFilteredPoints(elementIds []int, zwsTypeIDs []int, timestamp string) ([]Point, error) {
	zwsTypeIdsStr := strings.Trim(strings.Replace(fmt.Sprint(zwsTypeIDs), " ", ",", -1), "[]")

	url := fmt.Sprintf("%s/zulu/api/v1/filtered_points?zws_type_id=%v&timestamp=%v", c.serviceMapper.GetServiceURL("zulu"), zwsTypeIdsStr, timestamp)

	jsonPayload, err := json.Marshal(req{IDs: elementIds})
	if err != nil {
		return []Point{}, fmt.Errorf("ошибка маршалинга payload: %w", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))

	client := &http.Client{}
	client.Timeout = 100 * time.Second
	req.Header.Set("X-USER-ID", "1")
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
