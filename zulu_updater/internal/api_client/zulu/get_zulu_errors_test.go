package zulu

import (
	"context"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"zulu_updater/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestGetZuluErrors_Success(t *testing.T) {
	expectedResponse := models.ZWSErrorResponse{
		NetToolsTaskGetErrors: models.NetToolsTaskErrors{
			Errors: models.ErrorsSection{
				Count: 2,
				Errs: []models.Error{
					{Code: 123, Text: "First error"},
					{Code: 456, Text: "Second error"},
				},
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Contains(t, r.URL.Path, "/zws")
		w.WriteHeader(http.StatusOK)

		respBytes, _ := xml.Marshal(expectedResponse)
		w.Write(respBytes)
	}))
	defer ts.Close()

	os.Setenv("ZULU_BASE_URL", "http://"+ts.Listener.Addr().String())
	os.Setenv("ZULU_LAYER", "test-layer")
	os.Setenv("ZULU_TOKEN", "mockToken")

	client := NewApiClient()
	taskHandle := "test-task-handle"

	count, errors, err := client.GetZuluErrors(context.Background(), taskHandle)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.NetToolsTaskGetErrors.Errors.Count, count)
	assert.Equal(t, expectedResponse.NetToolsTaskGetErrors.Errors.Errs, errors)
}

func TestGetZuluErrors_RequestError(t *testing.T) {
	os.Setenv("ZULU_BASE_URL", "https://invalid-url")
	os.Setenv("ZULU_LAYER", "test-layer")
	os.Setenv("ZULU_TOKEN", "mockToken")

	client := NewApiClient()
	taskHandle := "test-task-handle"

	count, errors, err := client.GetZuluErrors(context.Background(), taskHandle)

	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.Nil(t, errors)
	assert.Contains(t, err.Error(), "Невозможно отправить запрос")
}

func TestGetZuluErrors_UnsuccessfulStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	os.Setenv("ZULU_BASE_URL", "http://"+ts.Listener.Addr().String())
	os.Setenv("ZULU_LAYER", "test-layer")
	os.Setenv("ZULU_TOKEN", "mockToken")

	client := NewApiClient()
	taskHandle := "test-task-handle"

	count, errors, err := client.GetZuluErrors(context.Background(), taskHandle)

	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.Nil(t, errors)
	assert.Contains(t, err.Error(), "код ответа")
}

func TestGetZuluErrors_InvalidXMLResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid-xml"))
	}))
	defer ts.Close()

	os.Setenv("ZULU_BASE_URL", "http://"+ts.Listener.Addr().String())
	os.Setenv("ZULU_LAYER", "test-layer")
	os.Setenv("ZULU_TOKEN", "mockToken")

	client := NewApiClient()
	taskHandle := "test-task-handle"

	count, errors, err := client.GetZuluErrors(context.Background(), taskHandle)

	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.Nil(t, errors)
	assert.Contains(t, err.Error(), "Ошибка парсинга XML ответа")
}
