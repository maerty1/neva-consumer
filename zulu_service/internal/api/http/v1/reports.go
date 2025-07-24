package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"zulu_service/internal/config/errors"
	reports_repository "zulu_service/internal/repositories/reports"
)

func RegisterReportsRouter(r *gin.Engine, reportsRepository reports_repository.Repository) {
	zulu := r.Group("/zulu/api/v1")
	{
		zulu.GET("/bi_dashboard/frame/1", func(ctx *gin.Context) {
			getBiDashboardFrame(ctx, reportsRepository)
		})
		zulu.GET("/bi_dashboard/frame/2/:element_id", func(ctx *gin.Context) {
			getBiDashboardFrameTwo(ctx, reportsRepository)
		})
		zulu.GET("/bi_dashboard/frame/3/:element_id", func(ctx *gin.Context) {
			getBiDashboardFrameThree(ctx, reportsRepository)
		})
		zulu.GET("/bi_dashboard/frame/3/:element_id/others", func(ctx *gin.Context) {
			getBiDashboardFrameThreeOthers(ctx, reportsRepository)
		})
	}

}

// @Router /zulu/api/v1/bi_dashboard/frame/1 [get]
// @Summary Получение первого фрейма BI дашборда
// @Description Каждый блок — это `Источник` (котельные)
// @Tags BI
// @Produce  json
// @Success 200 {object} []reports.BiDashboardFrame
func getBiDashboardFrame(ctx *gin.Context, reportsRepository reports_repository.Repository) {
	result, err := reportsRepository.GetBiDashboardFrame(ctx)
	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// @Router /zulu/api/v1/bi_dashboard/frame/2/{element_id} [get]
// @Summary Получение второго фрейма BI дашборда
// @Description Переменная element_id `Источник` получается из первого
// @Tags BI
// @Produce  json
// @Success 200 {object} []reports.BiDashboardFrame
func getBiDashboardFrameTwo(ctx *gin.Context, reportsRepository reports_repository.Repository) {

	elementID := ctx.Param("element_id")
	elementIDInt, err := strconv.Atoi(elementID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный elementID"})
		return
	}

	result, err := reportsRepository.GetBiDashboardFrameTwo(ctx, elementIDInt)
	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// @Router /zulu/api/v1/bi_dashboard/frame/3/{element_id} [get]
// @Summary Получение третьего фрейма BI дашборда
// @Description Переменная element_id `ЦТП` получается из второго кадра
// @Tags BI
// @Produce  json
// @Success 200 {object} []reports.BiDashboardFrame
func getBiDashboardFrameThree(ctx *gin.Context, reportsRepository reports_repository.Repository) {

	elementID := ctx.Param("element_id")
	elementIDInt, err := strconv.Atoi(elementID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный elementID"})
		return
	}

	result, err := reportsRepository.GetBiDashboardFrameThree(ctx, elementIDInt)
	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// @Router /zulu/api/v1/bi_dashboard/frame/3/{element_id}/others [get]
// @Summary Получение третьего фрейма BI дашборда при нажатии на "Остальное"
// @Description Переменная element_id `Источник` получается из первого кадра. Используется для детализации информации по блоку `Остальное` из второго фрейма
// @Tags BI
// @Produce  json
// @Success 200 {object} []reports.BiDashboardFrame
func getBiDashboardFrameThreeOthers(ctx *gin.Context, reportsRepository reports_repository.Repository) {

	elementID := ctx.Param("element_id")
	elementIDInt, err := strconv.Atoi(elementID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный elementID"})
		return
	}

	result, err := reportsRepository.GetBiDashboardFrameThreeOthers(ctx, elementIDInt)
	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}
	ctx.JSON(http.StatusOK, result)
}
