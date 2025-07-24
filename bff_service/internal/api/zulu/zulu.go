package zulu

import (
	"github.com/gin-gonic/gin"

	"bff_service/internal/facades/zulu"
)

func RegisterZuluRoutes(r *gin.Engine, zuluFacade zulu.Facade) {
	zulu := r.Group("/zulu/api/v1")

	{
		zulu.GET("/points", func(c *gin.Context) {
			zuluFacade.GetPoints(c)
		})
		zulu.GET("/points/with_rawdata", func(c *gin.Context) {
			zuluFacade.GetPointsWithRawdata(c)
		})
		zulu.GET("/points/:elem_id/full", func(c *gin.Context) {
			zuluFacade.GetPointsWithoutCopy(c)
		})
		zulu.GET("points/:elem_id/categories/:category_id", func(c *gin.Context) {
			zuluFacade.GetPointCategoryDataV2(c)
		})
	}
	zuluV2 := r.Group("/zulu/api/v2")
	{
		zuluV2.GET("points/categories/:category_id", func(c *gin.Context) {
			zuluFacade.GetPointsDataCategoryByZwsType(c)
		})
		zuluV2.GET("/points/:elem_id/full", func(c *gin.Context) {
			zuluFacade.GetFullPoint(c)
		})
		zuluV2.GET("points/:elem_id/categories/:category_id", func(c *gin.Context) {
			zuluFacade.GetPointCategoryDataV2(c)
		})
	}

}
