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

func TestPollMeasurePoints_RetryWithShortTimeout(t *testing.T) {
	attempts := 0

	// Создаем тестовый сервер, который имитирует API с 2 неудачными попытками и 1 успешной
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++

		var resp lers.PollMeasurePointsResponse

		// Первые две попытки возвращают PollSessionID = 0
		if attempts < 3 {
			resp = lers.PollMeasurePointsResponse{
				PollSessionID: 0,
				Status:        "Retry",
			}
		} else {
			resp = lers.PollMeasurePointsResponse{
				PollSessionID: 123,
				Status:        "Success",
			}
		}

		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}))
	defer ts.Close()

	client := lers.NewApiClient()

	timeout := time.Millisecond * 1
	response, err := client.PollMeasurePoints("test_token", ts.URL, []int{1, 2, 3}, "2024-01-01", "2024-01-02", timeout)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 123, response.PollSessionID)
	assert.Equal(t, 3, attempts)
}

func TestPollMeasurePoints_SuccessOnFirstTry(t *testing.T) {
	// Создаем тестовый сервер, который сразу возвращает успешный ответ
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := lers.PollMeasurePointsResponse{
			PollSessionID: 123,
			Status:        "Success",
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}))
	defer ts.Close()

	client := lers.NewApiClient()

	timeout := time.Millisecond * 1
	response, err := client.PollMeasurePoints("test_token", ts.URL, []int{1, 2, 3}, "2024-01-01", "2024-01-02", timeout)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 123, response.PollSessionID)
}

func TestPollMeasurePoints_FailureAfterRetries(t *testing.T) {
	// Создаем тестовый сервер, который всегда возвращает PollSessionID = 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := lers.PollMeasurePointsResponse{
			PollSessionID: 0,
			Status:        "Failure",
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}))
	defer ts.Close()

	client := lers.NewApiClient()

	timeout := time.Millisecond * 1
	response, err := client.PollMeasurePoints("test_token", ts.URL, []int{1, 2, 3}, "2024-01-01", "2024-01-02", timeout)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "PollSessionID == 0 после 7 попыток")
}

func TestPollMeasurePoints_UnexpectedStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := lers.NewApiClient()

	timeout := time.Millisecond * 1
	response, err := client.PollMeasurePoints("test_token", ts.URL, []int{1, 2, 3}, "2024-01-01", "2024-01-02", timeout)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "неожиданный status code")
}

func TestPollMeasurePoints_RequestCreationError(t *testing.T) {
	ts := httptest.NewServer(nil)
	defer ts.Close()

	client := lers.NewApiClient()

	// Вызов функции с некорректными данными для генерации ошибки создания запроса
	timeout := time.Millisecond * 1
	response, err := client.PollMeasurePoints("", string([]byte{0x7f}), []int{1, 2, 3}, "2024-01-01", "2024-01-02", timeout)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "ошибка создания запроса")
}
