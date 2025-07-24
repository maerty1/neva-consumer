package zulu

import (
	"golang.org/x/net/context"
	"net/http"
	"os"
	"zulu_updater/internal/models"
)

type ApiClient interface {
	UpdateAttribute(ctx context.Context, attributeName string, value float64, elemId int) error
	RunZuluCalculation(ctx context.Context) (string, error)
	GetZuluTaskStatus(ctx context.Context, taskHandle string) (string, error)
	GetZuluErrors(ctx context.Context, taskHandle string) (int, []models.Error, error)
	ExecuteSqlGetParametersByValNames(ctx context.Context, valNames []string, elemId int) (*models.Records, error)
}

var _ ApiClient = (*apiClient)(nil)

type apiClient struct {
	baseUrl string
	token   string
	layer   string
	client  *http.Client
}

func NewApiClient() ApiClient {
	var res apiClient
	res.baseUrl = os.Getenv("ZULU_BASE_URL")
	res.layer = os.Getenv("ZULU_LAYER")
	res.token = os.Getenv("ZULU_TOKEN")
	res.client = &http.Client{}
	return &res
}

func (a apiClient) setDefaultHeaders(req *http.Request) {
	headers := map[string]string{
		"Accept":           "*/*",
		"Authorization":    "Basic " + a.token,
		"Cache-Control":    "no-cache",
		"Content-Type":     "application/xml;charset=UTF-8",
		"X-Requested-With": "ZuluWebXMLHttpRequest",
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
}
