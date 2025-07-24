package app

import (
	"lers_integration_service/internal/api_clients/lers"
	"lers_integration_service/internal/db"
	"lers_integration_service/internal/repositories/measure_points"
	"lers_integration_service/internal/services/poller"
	"lers_integration_service/internal/services/retryer"
	"lers_integration_service/internal/services/synchronizer"
)

type serviceProvider struct {
	postgresDB              db.PostgresClient
	measurePointsRepository measure_points.Repository
	synchronizerService     synchronizer.Service
	pollerService           poller.Service
	retryerService          retryer.Service
	lersApiClient           lers.ApiClient
}

func newServiceProvider(db db.PostgresClient) *serviceProvider {
	return &serviceProvider{
		postgresDB: db,
	}
}

func (s *serviceProvider) PostgresDB() db.PostgresClient {
	return s.postgresDB
}

func (s *serviceProvider) MeasurePointsRepository() measure_points.Repository {
	if s.measurePointsRepository == nil {
		s.measurePointsRepository = measure_points.NewRepository(s.PostgresDB())
	}
	return s.measurePointsRepository
}

func (s *serviceProvider) LersApiClient() lers.ApiClient {
	if s.lersApiClient == nil {
		s.lersApiClient = lers.NewApiClient()
	}
	return s.lersApiClient
}

func (s *serviceProvider) SynchronizerService() synchronizer.Service {
	if s.synchronizerService == nil {
		s.synchronizerService = synchronizer.NewService(s.MeasurePointsRepository(), s.LersApiClient())
	}
	return s.synchronizerService
}

func (s *serviceProvider) PollerService() poller.Service {
	if s.pollerService == nil {
		s.pollerService = poller.NewService(s.MeasurePointsRepository(), s.LersApiClient())
	}
	return s.pollerService
}

func (s *serviceProvider) RetryerService() retryer.Service {
	if s.retryerService == nil {
		s.retryerService = retryer.NewService(s.MeasurePointsRepository(), s.LersApiClient())
	}
	return s.retryerService
}
