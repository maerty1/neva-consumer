package app

import (
	"log"
	"zulu_service/internal/config"
	"zulu_service/internal/db"
	geodataRepository "zulu_service/internal/repositories/geodata"
	reportsRepository "zulu_service/internal/repositories/reports"
	geodataService "zulu_service/internal/services/geodata"
)

type serviceProvider struct {
	postgresDB db.PostgresClient

	httpConfig config.HTTPConfig

	geodataRepository geodataRepository.Repository
	reportsRepository reportsRepository.Repository
	geodataService    geodataService.Service
}

func newServiceProvider(db db.PostgresClient) *serviceProvider {
	return &serviceProvider{
		postgresDB: db,
	}
}

func (s *serviceProvider) PostgresDB() db.PostgresClient {
	return s.postgresDB
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("не удалось получить конфигурацию http: %s", err.Error())
		}

		s.httpConfig = cfg
	}
	return s.httpConfig
}

func (s *serviceProvider) GeodataRepository() geodataRepository.Repository {
	if s.geodataRepository == nil {
		s.geodataRepository = geodataRepository.NewRepository(s.PostgresDB())
	}
	return s.geodataRepository
}

func (s *serviceProvider) ReportsRepository() reportsRepository.Repository {
	if s.reportsRepository == nil {
		s.reportsRepository = reportsRepository.NewRepository(s.PostgresDB())
	}
	return s.reportsRepository
}

func (s *serviceProvider) GeodataService() geodataService.Service {
	if s.geodataService == nil {
		s.geodataService = geodataService.NewService(s.GeodataRepository())
	}
	return s.geodataService
}
