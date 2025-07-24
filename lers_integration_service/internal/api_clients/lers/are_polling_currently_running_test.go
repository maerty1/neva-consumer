package lers_test

import (
	"encoding/json"
	"lers_integration_service/internal/api_clients/lers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArePollingsCurrentlyRunning_SuccessWithRunningPolls(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := []lers.Polling{
			{ID: 1},
			{ID: 2},
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}))
	defer ts.Close()

	client := lers.NewApiClient()

	running, err := client.ArePollingsCurrentlyRunning("test_token", ts.URL)

	assert.NoError(t, err)
	assert.True(t, running)
}

func TestArePollingsCurrentlyRunning_SuccessWithNoRunningPolls(t *testing.T) {
	// Создаем тестовый сервер, который возвращает пустой список задач опроса
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := []lers.Polling{}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}))
	defer ts.Close()

	client := lers.NewApiClient()

	running, err := client.ArePollingsCurrentlyRunning("test_token", ts.URL)

	assert.NoError(t, err)
	assert.False(t, running)
}

func TestArePollingsCurrentlyRunning_UnexpectedStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := lers.NewApiClient()

	running, err := client.ArePollingsCurrentlyRunning("test_token", ts.URL)

	assert.Error(t, err)
	assert.False(t, running)
	assert.Contains(t, err.Error(), "не удалось получить данные о задачах опроса")
}

func TestArePollingsCurrentlyRunning_RequestCreationError(t *testing.T) {
	// Создаем тестовый сервер, который никогда не будет вызван
	ts := httptest.NewServer(nil)
	defer ts.Close()

	client := lers.NewApiClient()

	// Вызов функции с некорректными данными для генерации ошибки создания запроса
	running, err := client.ArePollingsCurrentlyRunning("", string([]byte{0x7f}))

	assert.Error(t, err)
	assert.False(t, running)
	assert.Contains(t, err.Error(), "ошибка создания запроса")
}
