// Пакет poller предоставляет функциональность для выполнения опросов точек измерений (measure points) в ЛЭРС.
// Основная задача пакета — синхронизация данных с удалёнными серверами для различных аккаунтов.
// Пакет выполняет опрос точек измерений, управляет очередью запросов, обрабатывает результаты опросов и логирует успешные и неудачные попытки синхронизации в базу данных.
// Для выполнения этих операций используется взаимодействие с внешними API и репозиториями данных.
package poller

import (
	"context"
	"lers_integration_service/internal/api_clients/lers"
	"lers_integration_service/internal/repositories/measure_points"
)

type Service interface {
	RunPoller(ctx context.Context) error
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
