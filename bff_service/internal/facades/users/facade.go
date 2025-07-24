package users

import (
	"bff_service/internal/api_clients/users"
	"bff_service/internal/services"

	"github.com/gin-gonic/gin"
)

type Facade interface {
	Authenticate(c *gin.Context)
	GetAlisaToken(c *gin.Context)
}

var _ Facade = (*facade)(nil)

type facade struct {
	usersApiClient users.ApiClient
	jwtService     services.JwtService
}

func NewFacade(usersApiClient users.ApiClient, jwtService services.JwtService) Facade {
	return &facade{
		usersApiClient: usersApiClient,
		jwtService:     jwtService,
	}
}
