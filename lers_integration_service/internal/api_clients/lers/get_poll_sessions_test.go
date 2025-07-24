package lers_test

import (
	"encoding/json"
	"lers_integration_service/internal/api_clients/lers"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetPollSessions_AllSessionsCompleted(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := []map[string]interface{}{
			{
				"id":          1,
				"resultCode":  "Success",
				"endDateTime": "2024-01-01T00:00:00Z",
			},
			{
				"id":          2,
				"resultCode":  "Success",
				"endDateTime": "2024-01-01T01:00:00Z",
			},
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}))
	defer ts.Close()

	client := lers.NewApiClient()

	timeout := time.Millisecond * 1
	result, err := client.GetPollSessions("test_token", ts.URL, "2024-01-01", "2024-01-02", timeout)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Success", result[1])
	assert.Equal(t, "Success", result[2])
}

func TestGetPollSessions_UnfinishedSessions(t *testing.T) {
	attempts := 0

	// Создаем тестовый сервер, который сначала возвращает незавершенные сессии, а затем завершенные
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		var resp []map[string]interface{}

		if attempts < 2 {
			resp = []map[string]interface{}{
				{
					"id":         1,
					"resultCode": "InProgress",
				},
				{
					"id":         2,
					"resultCode": "InProgress",
				},
			}
		} else {
			resp = []map[string]interface{}{
				{
					"id":          1,
					"resultCode":  "Success",
					"endDateTime": "2024-01-01T00:00:00Z",
				},
				{
					"id":          2,
					"resultCode":  "Success",
					"endDateTime": "2024-01-01T01:00:00Z",
				},
			}
		}

		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}))
	defer ts.Close()

	client := lers.NewApiClient()

	timeout := time.Millisecond * 1
	result, err := client.GetPollSessions("test_token", ts.URL, "2024-01-01", "2024-01-02", timeout)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Success", result[1])
	assert.Equal(t, "Success", result[2])
	assert.Equal(t, 2, attempts)
}

func TestGetPollSessions_UnexpectedStatusCode(t *testing.T) {
	// Создаем тестовый сервер, который возвращает неожиданный статус код
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := lers.NewApiClient()

	timeout := time.Millisecond * 1
	result, err := client.GetPollSessions("test_token", ts.URL, "2024-01-01", "2024-01-02", timeout)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "неожиданный status code")
}

func TestGetPollSessions_RequestCreationError(t *testing.T) {
	// Создаем тестовый сервер, который никогда не будет вызван
	ts := httptest.NewServer(nil)
	defer ts.Close()

	client := lers.NewApiClient()

	// Вызов функции с некорректными данными для генерации ошибки создания запроса
	timeout := time.Millisecond * 1
	result, err := client.GetPollSessions("test_token", string([]byte{0x7f}), "2024-01-01", "2024-01-02", timeout)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "ошибка создания запроса")
}
