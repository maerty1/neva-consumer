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

func TestUpdateAttribute_Success(t *testing.T) {
	expectedResponse := models.ZWSUpdateAttributeResponse{
		UpdateElemAttributes: "Success",
		RetVal:               0,
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

	err := client.UpdateAttribute(context.Background(), "Temperature", 75.5, 123)

	assert.NoError(t, err)
}

func TestUpdateAttribute_RequestError(t *testing.T) {
	os.Setenv("ZULU_BASE_URL", "invalid-url")
	os.Setenv("ZULU_LAYER", "test-layer")
	os.Setenv("ZULU_TOKEN", "mockToken")

	client := NewApiClient()

	err := client.UpdateAttribute(context.Background(), "Temperature", 75.5, 123)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Невозможно отправить запрос")
}

func TestUpdateAttribute_UnsuccessfulStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	os.Setenv("ZULU_BASE_URL", "http://"+ts.Listener.Addr().String())
	os.Setenv("ZULU_LAYER", "test-layer")
	os.Setenv("ZULU_TOKEN", "mockToken")

	client := NewApiClient()

	err := client.UpdateAttribute(context.Background(), "Temperature", 75.5, 123)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "код ответа")
}

func TestUpdateAttribute_InvalidXMLResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid-xml"))
	}))
	defer ts.Close()

	os.Setenv("ZULU_BASE_URL", "http://"+ts.Listener.Addr().String())
	os.Setenv("ZULU_LAYER", "test-layer")
	os.Setenv("ZULU_TOKEN", "mockToken")

	client := NewApiClient()

	err := client.UpdateAttribute(context.Background(), "Temperature", 75.5, 123)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Ошибка парсинга XML ответа")
}

func TestUpdateAttribute_ReadBodyError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		hj, ok := w.(http.Hijacker)
		if ok {
			conn, _, _ := hj.Hijack()
			conn.Close()
		}
	}))
	defer ts.Close()

	os.Setenv("ZULU_BASE_URL", "http://"+ts.Listener.Addr().String())
	os.Setenv("ZULU_LAYER", "test-layer")
	os.Setenv("ZULU_TOKEN", "mockToken")

	client := NewApiClient()

	err := client.UpdateAttribute(context.Background(), "Temperature", 75.5, 123)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Невозможно прочесть ответ")
}
