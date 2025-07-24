package lers_test

import (
	"encoding/json"
	"lers_integration_service/internal/api_clients/lers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMeasurePoints_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := lers.MeasurePointsResponse{
			MeasurePoints: []lers.MeasurePoint{
				{ID: 1, Title: "Point 1"},
				{ID: 2, Title: "Point 2"},
			},
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}))
	defer ts.Close()

	client := lers.NewApiClient()

	measurePoints, err := client.GetMeasurePoints("test_token", ts.URL)

	assert.NoError(t, err)
	assert.NotNil(t, measurePoints)
	assert.Len(t, measurePoints, 2)
	assert.Equal(t, 1, measurePoints[0].ID)
	assert.Equal(t, "Point 1", measurePoints[0].Title)
	assert.Equal(t, 2, measurePoints[1].ID)
	assert.Equal(t, "Point 2", measurePoints[1].Title)
}

func TestGetMeasurePoints_UnexpectedStatusCode(t *testing.T) {
	// Создаем тестовый сервер, который возвращает неожиданный статус код
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := lers.NewApiClient()

	measurePoints, err := client.GetMeasurePoints("test_token", ts.URL)

	assert.Error(t, err)
	assert.Nil(t, measurePoints)
	assert.Contains(t, err.Error(), "не удалось получить точки измерения")
}

func TestGetMeasurePoints_InvalidJSONResponse(t *testing.T) {
	// Создаем тестовый сервер, который возвращает некорректный JSON
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"measurePoints":[{"id":1,"title":"Point 1"},{"id":`))
	}))
	defer ts.Close()

	client := lers.NewApiClient()

	measurePoints, err := client.GetMeasurePoints("test_token", ts.URL)

	assert.Error(t, err)
	assert.Nil(t, measurePoints)
	assert.Contains(t, err.Error(), "unexpected EOF")
}
func TestGetMeasurePoints_RequestCreationError(t *testing.T) {
	// Создаем тестовый сервер, который никогда не будет вызван
	ts := httptest.NewServer(nil)
	defer ts.Close()

	client := lers.NewApiClient()

	// Вызов функции с некорректными данными для генерации ошибки создания запроса
	measurePoints, err := client.GetMeasurePoints("", string([]byte{0x7f}))

	assert.Error(t, err)
	assert.Nil(t, measurePoints)
	assert.Contains(t, err.Error(), "invalid control character in URL")
}
