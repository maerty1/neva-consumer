package lers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PollMeasurePointsResponse struct {
	PollSessionID int    `json:"pollSessionId"`
	Status        string `json:"result"`
}

const retryCount = 7
const Timeout = 10 * time.Second

// PollMeasurePoints создает запрос на опрос точек измерения за период.
// Если возникает ошибка, то делается ретрай
func (c *apiClient) PollMeasurePoints(token string, serverHost string, measurePointIDs []int, startDate, endDate string, timeout time.Duration) (*PollMeasurePointsResponse, error) {
	url := fmt.Sprintf("%s/api/v0.1/Poll/ManualPoll/Archive/%s/%s", serverHost, startDate, endDate)

	payload := map[string]interface{}{
		"requestedDataTypes": []string{"Hour"},
		"absentDataOnly":     true,
		"additional": map[string]interface{}{
			"detectDevice":          false,
			"connectionTimeOut":     600,
			"debugEnabled":          false,
			"gprsAutoDisconnect":    false,
			"performTimeCorrection": false,
			"ignoreTimeDifference":  false,
			"responseDelay":         0,
			"controlLoad":           "None",
		},
		"measurePoints":    measurePointIDs,
		"userConnectionId": "",
		"startType":        "Normal",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга payload: %w", err)
	}

	for attempt := 1; attempt <= retryCount; attempt++ {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			return nil, fmt.Errorf("ошибка создания запроса: %w", err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("ошибка отправки запроса: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("неожиданный status code: %d", resp.StatusCode)
		}

		var response PollMeasurePointsResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
		}

		// Проверяем, если PollSessionID == 0, делаем ретрай
		if response.PollSessionID == 0 {
			if attempt < retryCount {
				fmt.Printf("PollSessionID == 0, повтор попытки %d через 10 секунд...\n", attempt)
				time.Sleep(timeout)
				continue
			} else {
				return nil, fmt.Errorf("PollSessionID == 0 после 7 попыток")
			}
		}

		// Если PollSessionID не равен 0, возвращаем результат
		return &response, nil
	}

	// Этот код не должен быть достигнут, но на всякий случай
	return nil, fmt.Errorf("не удалось получить корректный PollSessionID после 7 попыток")
}

// {
//     "requestedDataTypes": [
//         "Hour"
//     ],
//     "absentDataOnly": true,
//     "additional": {
//         "detectDevice": false,
//         "connectionTimeOut": 75,
//         "debugEnabled": false,
//         "gsmProtocol": "ModemDefined",
//         "gprsAutoDisconnect": false,
//         "performTimeCorrection": false,
//         "ignoreTimeDifference": false,
//         "responseDelay": 0
//     },
//     "measurePoints": [
//         743
//     ],
//     "userConnectionId": "",
//     "startMode": "Normal",
//     "pollConnectionId": 613
// }
