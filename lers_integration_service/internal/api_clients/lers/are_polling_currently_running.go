package lers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Task struct {
	Type string `json:"type"`
}

type Polling struct {
	ID   int  `json:"id"`
	Task Task `json:"task"`
}

// ArePollingsCurrentlyRunning возвращает статус, находятся ли точки опроса (measure points) в процессе получения данных от удаленного сервера
func (c *apiClient) ArePollingsCurrentlyRunning(token string, serverHost string) (bool, error) {
	url := fmt.Sprintf("%s/api/v0.1/Poll/PollTasks/View", serverHost)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("не удалось получить данные о задачах опроса: %s", resp.Status)
	}

	var response []Polling
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	// Проходим по всем задачам и проверяем их тип
	for _, polling := range response {
		// Если тип задачи не "Auto", считаем, что опросы запущены
		if polling.Task.Type != "Auto" {
			return true, nil
		}
	}

	// Если все задачи типа "Auto" или список пуст, считаем, что очередь на опрос пуста
	return false, nil
}
