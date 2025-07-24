package zulu

import (
	"context"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"zulu_updater/internal/models"
)

func TestExecuteSqlGetParametersByValNames_Success(t *testing.T) {
	expectedRecords := models.Records{
		Record: []models.Record{
			{Field: []models.Field{{"name", "val"}}},
		},
	}
	responseXML := models.ZwsSqlResponse{
		LayerExecSql: struct {
			Records models.Records `xml:"Records"`
		}{
			Records: expectedRecords,
		},
	}

	responseBody, _ := xml.Marshal(responseXML)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/zws")
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		assert.Contains(t, string(body), "<LayerExecSql>")
		assert.Contains(t, string(body), "SELECT val1,val2 where Sys=123")

		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)
	}))
	defer ts.Close()

	os.Setenv("ZULU_BASE_URL", ts.URL)
	os.Setenv("ZULU_LAYER", "test-layer")
	os.Setenv("ZULU_TOKEN", "mockToken")

	client := NewApiClient()

	valNames := []string{"val1", "val2"}
	elemId := 123
	ctx := context.Background()

	result, err := client.ExecuteSqlGetParametersByValNames(ctx, valNames, elemId)

	assert.NoError(t, err)
	assert.Equal(t, &expectedRecords, result)
}

func TestExecuteSqlGetParametersByValNames_RequestError(t *testing.T) {
	os.Setenv("ZULU_BASE_URL", "invalid-url")
	os.Setenv("ZULU_LAYER", "test-layer")
	os.Setenv("ZULU_TOKEN", "mockToken")

	client := NewApiClient()

	valNames := []string{"val1", "val2"}
	elemId := 123
	ctx := context.Background()

	result, err := client.ExecuteSqlGetParametersByValNames(ctx, valNames, elemId)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Невозможно отправить запрос")
}

func TestExecuteSqlGetParametersByValNames_UnsuccessfulStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	os.Setenv("ZULU_BASE_URL", ts.URL)
	os.Setenv("ZULU_LAYER", "test-layer")
	os.Setenv("ZULU_TOKEN", "mockToken")

	client := NewApiClient()

	valNames := []string{"val1", "val2"}
	elemId := 123
	ctx := context.Background()

	result, err := client.ExecuteSqlGetParametersByValNames(ctx, valNames, elemId)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "код ответа")
}

func TestExecuteSqlGetParametersByValNames_InvalidXMLResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid-xml"))
	}))
	defer ts.Close()

	os.Setenv("ZULU_BASE_URL", ts.URL)
	os.Setenv("ZULU_LAYER", "test-layer")
	os.Setenv("ZULU_TOKEN", "mockToken")

	client := NewApiClient()

	valNames := []string{"val1", "val2"}
	elemId := 123
	ctx := context.Background()

	result, err := client.ExecuteSqlGetParametersByValNames(ctx, valNames, elemId)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Ошибка парсинга XML ответа")
}
