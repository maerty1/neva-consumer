package config

import (
	"os"
)

const UsersHTTPServerUrl = "USERS_HTTP_SERVICE_URL"
const CoreDataHTTPServerUrl = "CORE_DATA_SERVICE_HTTP_SERVICE_URL"
const ZuluHTTPServerUrl = "ZULU_SERVICE_HTTP_SERVICE_URL"

type ServiceMapper interface {
	GetServiceURL(service string) string
}

type serviceMapper struct {
	Users    string
	CoreData string
	Zulu     string
}

func NewServiceMapper() ServiceMapper {
	return &serviceMapper{
		Users:    getEnv(UsersHTTPServerUrl, "http://host.docker.internal:8001"),
		CoreData: getEnv(CoreDataHTTPServerUrl, "http://host.docker.internal:8002"),
		Zulu:     getEnv(ZuluHTTPServerUrl, "http://host.docker.internal:8004"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func (sm *serviceMapper) GetServiceURL(service string) string {
	switch service {
	case "users":
		return sm.Users
	case "core":
		return sm.CoreData
	case "zulu":
		return sm.Zulu
	default:
		return ""
	}
}
