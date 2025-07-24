package lers_test

import (
	"encoding/json"
	"lers_integration_service/internal/api_clients/lers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConsumptionData_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := lers.ConsumptionResponse{
			PressureType: "High",
			HourConsumption: []lers.ConsumptionData{
				{
					ResourceKind: "Electricity",
					DateTime:     "2024-08-27T00:00:00Z",
					Values: []lers.Value{
						{
							DataParameter: "ActiveEnergy",
							Value:         123.45,
							IsBad:         false,
							IsCalc:        false,
						},
					},
				},
			},
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}))
	defer ts.Close()

	client := lers.NewApiClient()

	response, err := client.GetConsumptionData(1, "test_token", ts.URL, 123, "2024-01-01", "2024-01-02")

	data := response.HourConsumption

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Len(t, data, 1)
	assert.Equal(t, "Electricity", data[0].ResourceKind)
	assert.Equal(t, "2024-08-27T00:00:00Z", data[0].DateTime)
	assert.Len(t, data[0].Values, 1)
	assert.Equal(t, "ActiveEnergy", data[0].Values[0].DataParameter)
	assert.Equal(t, 123.45, data[0].Values[0].Value)
}

func TestGetConsumptionData_UnexpectedStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := lers.NewApiClient()

	data, err := client.GetConsumptionData(1, "test_token", ts.URL, 123, "2024-01-01", "2024-01-02")

	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "не удалось получить данные о потреблении")
}

func TestGetConsumptionData_RequestCreationError(t *testing.T) {
	ts := httptest.NewServer(nil)
	defer ts.Close()

	client := lers.NewApiClient()

	// Вызов функции с некорректными данными для генерации ошибки создания запроса
	data, err := client.GetConsumptionData(1, "test_token", string([]byte{0x7f}), 123, "2024-01-01", "2024-01-02")

	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestGetConsumptionData_InvalidJSONResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"invalidJson":`))
	}))
	defer ts.Close()

	client := lers.NewApiClient()

	data, err := client.GetConsumptionData(1, "test_token", ts.URL, 123, "2024-01-01", "2024-01-02")

	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "unexpected EOF")
}
