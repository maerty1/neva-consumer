package lers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// GetPollSessions возвращает историю опросов точек измерения.
// Если какие-то из опросов находятся в процессе получения данных, функция дожидается завершения этого процесса
func (c *apiClient) GetPollSessions(token string, serverHost string, startDate string, endDate string, timeout time.Duration) (map[int]string, error) {
	url := fmt.Sprintf("%s/api/v0.1/Poll/PollSessions/%s/%s", serverHost, startDate, endDate)

	for {
		req, err := http.NewRequest("GET", url, nil)
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

		// Используем слайс для временного хранения ответа
		var response []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
		}

		// Проверяем, есть ли незавершенные сессии
		allSessionsCompleted := true
		for _, session := range response {
			if session["endDateTime"] == nil {
				allSessionsCompleted = false
				break
			}
		}

		// Если все сессии завершены, формируем результат и возвращаем
		if allSessionsCompleted {
			resultMap := make(map[int]string)
			for _, session := range response {
				pollSessionID := int(session["id"].(float64))
				resultCode := session["resultCode"].(string)
				resultMap[pollSessionID] = resultCode
			}
			return resultMap, nil
		}

		// Если есть незавершенные сессии, ждем и повторяем запрос
		log.Println("Незавершенные сессии обнаружены, ожидание 30 секунд перед повторным запросом...")
		time.Sleep(timeout)
	}
}
