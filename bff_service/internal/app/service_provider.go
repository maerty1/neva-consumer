package app

import (
	coreDataApiClient "bff_service/internal/api_clients/core"
	usersApiClient "bff_service/internal/api_clients/users"
	zuluApiClient "bff_service/internal/api_clients/zulu"
	"bff_service/internal/config"
	"bff_service/internal/facades/core"
	"bff_service/internal/facades/users"
	"bff_service/internal/facades/zulu"
	"bff_service/internal/handlers/bff"
	"bff_service/internal/services"
	"log"
	"time"
)

type serviceProvider struct {
	httpConfig          config.HTTPConfig
	jwtConfig           config.JWTConfig
	serviceMapperConfig config.ServiceMapper

	bffHandler bff.BFFHandler

	usersFacade users.Facade
	coreFacade  core.Facade
	zuluFacade  zulu.Facade

	usersApiClient          usersApiClient.ApiClient
	zuluApiClient           zuluApiClient.ApiClient
	cachedZuluApiClient     zuluApiClient.ApiClient
	coreDataApiClient       coreDataApiClient.ApiClient
	cachedCoreDataApiClient coreDataApiClient.ApiClient

	jwtService services.JwtService
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
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
func (s *serviceProvider) JWTConfig() config.JWTConfig {
	if s.jwtConfig == nil {
		cfg, err := config.NewJWTConfig()
		if err != nil {
			log.Fatalf("не удалось получить конфигурацию jwt: %s", err.Error())
		}

		s.jwtConfig = cfg
	}
	return s.jwtConfig
}

func (s *serviceProvider) BFFHandler() bff.BFFHandler {
	if s.bffHandler == nil {
		s.bffHandler = bff.NewBFFHandlerExperimental(3, 2*time.Second, s.UsersApiClient())
	}

	return s.bffHandler
}

func (s *serviceProvider) UsersFacade() users.Facade {
	if s.usersFacade == nil {
		s.usersFacade = users.NewFacade(s.UsersApiClient(), s.JWTService())
	}

	return s.usersFacade
}

func (s *serviceProvider) CoreFacade() core.Facade {
	if s.coreFacade == nil {
		s.coreFacade = core.NewFacade()
	}

	return s.coreFacade
}

func (s *serviceProvider) ZuluFacade() zulu.Facade {
	if s.zuluFacade == nil {
		s.zuluFacade = zulu.NewFacade(s.CachedZuluApiClient(), s.CachedCoreApiClient())
	}

	return s.zuluFacade
}

func (s *serviceProvider) CachedZuluApiClient() zuluApiClient.ApiClient {
	if s.cachedZuluApiClient == nil {
		baseZuluClient := s.ZuluApiClient()
		s.cachedZuluApiClient = zuluApiClient.NewCacheApiClient(
			baseZuluClient,
			5*time.Minute,
		)
	}

	return s.cachedZuluApiClient
}
func (s *serviceProvider) CachedCoreApiClient() coreDataApiClient.ApiClient {
	if s.cachedCoreDataApiClient == nil {
		baseCoreClient := s.CoreDataApiClient()
		s.cachedCoreDataApiClient = coreDataApiClient.NewCacheApiClient(
			baseCoreClient,
			5*time.Minute,
		)
	}

	return s.cachedCoreDataApiClient
}

func (s *serviceProvider) UsersApiClient() usersApiClient.ApiClient {
	if s.usersApiClient == nil {
		s.usersApiClient = usersApiClient.NewApiClient(s.ServiceMapper())
	}

	return s.usersApiClient
}

func (s *serviceProvider) CoreDataApiClient() coreDataApiClient.ApiClient {
	if s.coreDataApiClient == nil {
		s.coreDataApiClient = coreDataApiClient.NewApiClient(s.ServiceMapper())
	}

	return s.coreDataApiClient
}
func (s *serviceProvider) ZuluApiClient() zuluApiClient.ApiClient {
	if s.zuluApiClient == nil {
		s.zuluApiClient = zuluApiClient.NewApiClient(s.ServiceMapper())
	}

	return s.zuluApiClient
}

func (s *serviceProvider) ServiceMapper() config.ServiceMapper {
	if s.serviceMapperConfig == nil {
		s.serviceMapperConfig = config.NewServiceMapper()
	}

	return s.serviceMapperConfig
}

func (s *serviceProvider) JWTService() services.JwtService {
	if s.jwtService == nil {
		s.jwtService = services.NewService(s.JWTConfig())
	}

	return s.jwtService
}
