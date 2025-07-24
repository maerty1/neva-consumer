// Пакет retryer предоставляет функциональность для повторного выполнения опросов точек измерений (measure_points) в ЛЭРС.
// Он отвечает за обработку опросов (polls), требующих повторного опроса, управление процессом повторных запросов, а также обновление статусов опросов.
// Пакет использует взаимодействие с внешними API и репозиториями данных для выполнения ретраев и логирования результатов в базу данных.
package retryer

import (
	"context"
	"lers_integration_service/internal/api_clients/lers"
	"lers_integration_service/internal/repositories/measure_points"
)

type Service interface {
	RunRetryer(ctx context.Context) error
}

var _ Service = (*service)(nil)

type service struct {
	measurePointsRepository measure_points.Repository
	lersApiClient           lers.ApiClient
}

func NewService(repository measure_points.Repository, apiClient lers.ApiClient) Service {
	return &service{
		measurePointsRepository: repository,
		lersApiClient:           apiClient,
	}
}
