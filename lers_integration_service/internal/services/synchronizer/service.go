// Пакет synchronizer предназначен для синхронизации данных о потреблении ресурсов с серверов ЛЭРС.
// Основная задача пакета — получение и сохранение данных о потреблении для каждой точки измерения (measure point), а также  логирование успешных и неудачных попыток синхронизации.
// Пакет взаимодействует с внешними API для получения данных и использует репозитории для их сохранения в базу данных.
package synchronizer

import (
	"context"
	"lers_integration_service/internal/api_clients/lers"
	"lers_integration_service/internal/repositories/measure_points"
)

type Service interface {
	RunSynchronizer(ctx context.Context) error
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
