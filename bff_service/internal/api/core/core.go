package core

import (
	"bff_service/internal/facades/core"

	"github.com/gin-gonic/gin"
)

func RegisterCoreRoutes(r *gin.Engine, coreFacade core.Facade) {
	core := r.Group("/core/api/v1")

	{

		core.GET("/weather/current", func(c *gin.Context) {
			coreFacade.GetCurrentWeather(c)
		})

	}
}
