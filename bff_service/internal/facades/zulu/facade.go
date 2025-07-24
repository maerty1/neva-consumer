package zulu

import (
	"github.com/gin-gonic/gin"

	"bff_service/internal/api_clients/core"
	"bff_service/internal/api_clients/zulu"
)

type Facade interface {
	GetPoints(c *gin.Context)
	GetFullPoint(c *gin.Context)
	GetPointsWithoutCopy(c *gin.Context)
	GetPointCategoryData(c *gin.Context)
	GetPointsWithRawdata(c *gin.Context)
	GetPointsDataCategoryByZwsType(c *gin.Context)

	GetPointCategoryDataV2(c *gin.Context)
}

var _ Facade = (*facade)(nil)

type facade struct {
	zuluApiClient     zulu.ApiClient
	coreDataApiClient core.ApiClient
}

func NewFacade(zuluApiClient zulu.ApiClient, coreDataApiClient core.ApiClient) Facade {
	return &facade{
		zuluApiClient:     zuluApiClient,
		coreDataApiClient: coreDataApiClient,
	}
}
