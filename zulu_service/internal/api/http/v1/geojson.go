package v1

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"zulu_service/internal/config/errors"
	"zulu_service/internal/models/geodata"
	geodata_repository "zulu_service/internal/repositories/geodata"
)

func RegisterGeoJsonRouter(r *gin.Engine, geodataRepository geodata_repository.Repository) {
	zulu := r.Group("/zulu/api/v1")
	{
		zulu.GET("/geojson", func(ctx *gin.Context) {
			getGeoJson(ctx, geodataRepository)
		})
		zulu.GET("/pipeline_depths", func(ctx *gin.Context) {
			getPipelineDepths(ctx, geodataRepository)
		})
		zulu.GET("/states", func(ctx *gin.Context) {
			getStates(ctx, geodataRepository)
		})
		zulu.GET("/points", func(ctx *gin.Context) {
			getPoints(ctx, geodataRepository)
		})
		zulu.POST("/filtered_points", func(ctx *gin.Context) {
			getFilteredPoints(ctx, geodataRepository)
		})
		zulu.GET("/points/:element_id/full", func(ctx *gin.Context) {
			getPointFull(ctx, geodataRepository)
		})
		zulu.GET("/elements/:element_id", func(ctx *gin.Context) {
			getElementData(ctx, geodataRepository)
		})
		zulu.GET("/points/:element_id/categories", func(ctx *gin.Context) {
			getPointCategories(ctx, geodataRepository)
		})
		zulu.GET("/points/:element_id/categories/:category_id", func(ctx *gin.Context) {
			getPointDataByCategory(ctx, geodataRepository)
		})
		zulu.GET("/points/:element_id/icon", func(ctx *gin.Context) {
			getIcon(ctx, geodataRepository)
		})
		zulu.GET("/points/:element_id/schema", func(ctx *gin.Context) {
			getSchema(ctx, geodataRepository)
		})
		zulu.GET("/enums/measurement_groups", func(ctx *gin.Context) {
			getMeasurementGroups(ctx, geodataRepository)
		})

		zulu.POST("/points/categories/:category_id", func(ctx *gin.Context) {
			getPointsCategory(ctx, geodataRepository)
		})
		zulu.GET("/points/categories/:category_id", func(ctx *gin.Context) {
			getPointsCategoryByZwsType(ctx, geodataRepository)
		})
	}
	zuluV2 := r.Group("/zulu/api/v2")
	{
		zuluV2.GET("/geojson", func(ctx *gin.Context) {
			getGeoJsonV2(ctx, geodataRepository)
		})
		zuluV2.GET("/points/categories/:category_id", func(ctx *gin.Context) {
			getPointsCategoryByZwsTypeV2(ctx, geodataRepository)
		})
	}

}

// @Router /zulu/api/v1/geojson [get]
// @Summary Получение GeoJSON данных
// @Description Возвращает коллекцию географических объектов в формате GeoJSON.
// @Tags Geojson
// @Produce  json
// @Success 200 {object} geojson.GeoJSONFeatureCollection "GeoJSON FeatureCollection"
func getGeoJson(ctx *gin.Context, geodataRepository geodata_repository.Repository) {
	geoJSON, err := geodataRepository.GetGeoJson(ctx)
	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}
	ctx.Data(http.StatusOK, "application/json", geoJSON)
}

// @Router /zulu/api/v1/pipeline_depths [get]
// @Summary Получение глубину нахождения труб под землей
// @Description Получение глубину нахождения труб под землей
// @Tags Pipeline
// @Produce  json
// @Success 200 {object} map[int]geodata.PipelineDepth
func getPipelineDepths(ctx *gin.Context, geodataRepository geodata_repository.Repository) {
	pipelineDepth, err := geodataRepository.GetPipelineDepth(ctx)
	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}
	ctx.JSON(http.StatusOK, pipelineDepth)
}

// @Router /zulu/api/v2/geojson [get]
// @Summary Получение GeoJSON данных
// @Description Возвращает коллекцию географических объектов в формате GeoJSON.
// @Tags Geojson
// @Produce  json
// @Success 200 {object} geojson.GeoJSONFeatureCollectionV2 "GeoJSON FeatureCollection"
func getGeoJsonV2(ctx *gin.Context, geodataRepository geodata_repository.Repository) {
	geoJSON, err := geodataRepository.GetGeoJsonV2(ctx)
	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}
	ctx.Data(http.StatusOK, "application/json", geoJSON)
}

// @Router /zulu/api/v1/states [get]
// @Summary Получение состояний объектов
// @Description Возвращает список состояний объектов из словаря `zulu.dict_object_states`.
// @Tags Geojson
// @Produce  json
// @Success 200 {object} geodata.ObjectStatesResponse "Список состояний объектов"
func getStates(ctx *gin.Context, geodataRepository geodata_repository.Repository) {
	states, err := geodataRepository.GetStates(ctx)
	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}

	response := geodata.ObjectStatesResponse{
		States: states,
	}

	ctx.JSON(http.StatusOK, response)
}

// @Router /zulu/api/v1/elements/{element_id} [get]
// @Summary Получение значений объектов
// @Tags Geojson
// @Produce  json
// @Success 200 {object} geodata.ElementData "Список значений объектов"
func getElementData(ctx *gin.Context, geodataRepository geodata_repository.Repository) {
	elementID := ctx.Param("element_id")
	elementIDInt, err := strconv.Atoi(elementID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный elementID"})
		return
	}

	data, err := geodataRepository.GetElementDataByID(ctx, elementIDInt)

	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

// @Router /zulu/api/v1/points [get]
// @Summary Получение значений объектов
// @Tags Points
// @Param zws_type_id query []int false "ID типов объектов" collectionFormat(csv)
// @Produce  json
// @Success 200 {object} []geodata.Point
func getPoints(ctx *gin.Context, geodataRepository geodata_repository.Repository) {
	zwsTypeIdsStr := ctx.Query("zws_type_id")

	var ids []int
	if zwsTypeIdsStr == "" {
		ids = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	} else {
		for _, idStr := range strings.Split(zwsTypeIdsStr, ",") {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "неверны zws_type_id. Доступно: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11"})
				return
			}
			ids = append(ids, id)
		}
	}

	data, err := geodataRepository.GetPoints(ctx, ids)

	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

// @Router /zulu/api/v1/filtered_points [post]
// @Summary Получение значений объектов
// @Param zws_type_id query []int false "ID типов объектов для фильтрации" collectionFormat(csv)
// @Tags Internal
// @Produce  json
// @Success 200 {object} []geodata.Point
func getFilteredPoints(ctx *gin.Context, geodataRepository geodata_repository.Repository) {
	var request struct {
		IDs []int `json:"ids"`
	}

	// Попытка привязать JSON тело запроса к структуре
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	var ids []int
	zwsTypeIdsStr := ctx.Query("zws_type_id")

	if zwsTypeIdsStr == "" {
		ids = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	} else {
		for _, idStr := range strings.Split(zwsTypeIdsStr, ",") {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "неверны zws_type_id. Доступно: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11"})
				return
			}
			ids = append(ids, id)
		}
	}
	timestamp := ctx.Query("timestamp")

	data, err := geodataRepository.GetFilteredPoints(ctx, request.IDs, ids, timestamp)

	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

// @Router /zulu/api/v1/points/{element_id}/full [get]
// @Summary ...
// @Tags Internal
// @Produce  json
// @Success 200 {object} geodata.FullElementData
func getPointFull(ctx *gin.Context, geodataRepository geodata_repository.Repository) {
	elementID := ctx.Param("element_id")
	elementIDInt, err := strconv.Atoi(elementID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный elementID"})
		return
	}

	nDays := ctx.Query("n_days")
	fmt.Println(nDays)

	nDaysInt, err := strconv.Atoi(elementID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный n_days"})
		return
	}

	// TODO: Убрать хардкод
	data, err := geodataRepository.GetFullByElemID(ctx, elementIDInt, nDaysInt)

	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

// @Router /zulu/api/v1/enums/measurement_groups [get]
// @Summary Получение значений объектов
// @Description Возвращает словарь групп измерений, где ключ — ID группы. `{"1": {"name": "Температура", "unit": "°C"}, "2": {"name": "Давление", "unit": "атм"}}`
// @Tags Enums
// @Produce  json
// @Success 200 {object} map[int]geodata.MeasurementGroupEnum
func getMeasurementGroups(ctx *gin.Context, geodataRepository geodata_repository.Repository) {
	data, err := geodataRepository.GetMeasurementGroupsEnum(ctx)
	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

// @Router /zulu/api/v1/points/{element_id}/categories [get]
// @Tags Categories
// @Produce  json
// @Success 200 {object} geodata.PointWithCategories
func getPointCategories(ctx *gin.Context, geodataRepository geodata_repository.Repository) {
	elementID := ctx.Param("element_id")
	elementIDInt, err := strconv.Atoi(elementID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный elementID"})
		return
	}
	data, err := geodataRepository.GetPointCategories(ctx, elementIDInt)

	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

// @Router /zulu/api/v1/points/{element_id}/categories/{category_id} [get]
// @Tags Internal
// @Produce  json
// @Param type query string true "Type of data" Enums(group, keyvalue)
// @Success 200 {object} geodata.GetPointDataByCategoryGroup "Ответ для type=group"
// @Success 201 {object} geodata.GetPointDataByCategoryKeyvalue "Ответ для type=keyvalue"
func getPointDataByCategory(ctx *gin.Context, geodataRepository geodata_repository.Repository) {
	elementID := ctx.Param("element_id")
	elementIDInt, err := strconv.Atoi(elementID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный elementID"})
		return
	}

	categoryID := ctx.Param("category_id")
	categoryIDInt, err := strconv.Atoi(categoryID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный categoryID"})
		return
	}

	timestamp := ctx.Query("timestamp")
	nDays := ctx.Query("n_days")

	// Проверка наличия только одного из параметров
	if (timestamp == "" && nDays == "") || (timestamp != "" && nDays != "") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Необходимо указать либо 'timestamp', либо 'n_days', но не оба одновременно.",
		})
		return
	}

	var nDaysInt int
	if nDays != "" {
		nDaysInt, err = strconv.Atoi(nDays)
		if err != nil || nDaysInt < 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный nDays. Должно быть целое число больше 0."})
			return
		}
	}

	categroyType := ctx.Query("type")
	if categroyType == "group" {
		data, err := geodataRepository.GetPointDataByCategoryGroup(ctx, elementIDInt, categoryIDInt, timestamp, nDaysInt)
		ctx.JSON(http.StatusOK, data)

		if err != nil {
			resp := errors.GetHTTPStatus(err)
			ctx.JSON(resp.Status, gin.H{"error": resp.Message})
			return
		}

	} else if categroyType == "keyvalue" {
		data, err := geodataRepository.GetPointDataByCategoryKeyvalue(ctx, elementIDInt, categoryIDInt)
		ctx.JSON(http.StatusOK, data)

		if err != nil {
			resp := errors.GetHTTPStatus(err)
			ctx.JSON(resp.Status, gin.H{"error": resp.Message})
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный queryParam type. Доступно: keyvalue, group"})
		return
	}
}

// @Router /zulu/api/v1/points/{element_id}/icon [get]
// @Summary      Получить иконку по элементу
// @Description  Возвращает файл иконки для заданного элемента (element_id)
// @Tags         Schemas
// @Param        element_id   path      int  true  "ID элемента"
// @Success      200          {file}    file  "Файл иконки"
// @Failure      400          {object}  map[string]string  "Некорректный elementID"
// @Failure      404          {object}  map[string]string  "Файл не найден"
func getIcon(ctx *gin.Context, geodataRepository geodata_repository.Repository) {
	elementID := ctx.Param("element_id")
	elementIDInt, err := strconv.Atoi(elementID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный elementID"})
		return
	}

	iconID, err := geodataRepository.GetIconIdByZwsType(ctx, elementIDInt)
	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}

	basePath := "/code/files/icons"
	extensions := []string{".jpg", ".png", ".jpeg"}

	id := strconv.Itoa(iconID)

	var filePath string
	for _, ext := range extensions {
		path := filepath.Join(basePath, id+ext)
		if _, err := os.Stat(path); err == nil {
			filePath = path
			break
		}
	}

	if filePath == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	ctx.File(filePath)
}

// @Router /zulu/api/v1/points/{element_id}/schema [get]
// @Summary      Получить схему по элементу
// @Description  Возвращает файл схемы для заданного элемента (element_id)
// @Tags         Schemas
// @Param        element_id   path      int  true  "ID элемента"
// @Param        theme query string Enums(light, dark)
// @Success      200          {file}    file  "Файл схемы"
// @Failure      400          {object}  map[string]string  "Некорректный elementID"
// @Failure      404          {object}  map[string]string  "Файл не найден"
func getSchema(ctx *gin.Context, geodataRepository geodata_repository.Repository) {
	// Времянка пока нет s3

	elementID := ctx.Param("element_id")
	elementIDInt, err := strconv.Atoi(elementID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный elementID"})
		return
	}

	schemaID, err := geodataRepository.GetSchemaIdByElemId(ctx, elementIDInt)
	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}

	theme := ctx.Query("theme")

	if theme == "" {
		theme = "light/"
	} else {
		theme = "dark/"
	}

	basePath := "/code/files/schemas/" + theme
	extensions := []string{".jpg", ".png", ".jpeg", ".svg"}

	id := strconv.Itoa(schemaID)

	var filePath string
	for _, ext := range extensions {
		path := filepath.Join(basePath, id+ext)

		fmt.Println(path)

		if _, err := os.Stat(path); err == nil {
			filePath = path
			break
		}
	}

	if filePath == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	ctx.File(filePath)
}

// @Router /zulu/api/v1/points/categories/:category_id [post]
// @Summary ...
// @Tags Internal
// @Produce  json
// @Success 200 {object} map[int]geodata.GetPointsDataByCategoryGroup
func getPointsCategory(ctx *gin.Context, geodataRepository geodata_repository.Repository) {

	categoryID := ctx.Param("category_id")
	categoryIDInt, err := strconv.Atoi(categoryID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный categoryID"})
		return
	}

	var elemIDs []int

	if err := ctx.BindJSON(&elemIDs); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные elemIDs"})
		return
	}

	data, err := geodataRepository.GetPointsDataByCategoryGroup(ctx, elemIDs, categoryIDInt)

	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

// @Router /zulu/api/v1/points/categories/:category_id [get]
// @Summary Получение значений объектов
// @Tags Internal
// @Param zws_type_id query []int false "ID типов объектов для фильтрации" collectionFormat(csv)
// @Param timestamp query string false "Timestamp для получения данных за конкретный период"
// @Produce  json
// @Success 200 {object} map[int]geodata.GetPointsDataByCategoryGroup
func getPointsCategoryByZwsType(ctx *gin.Context, geodataRepository geodata_repository.Repository) {
	zwsTypeIdsStr := ctx.Query("zws_type_id")

	var zwsIDs []int
	if zwsTypeIdsStr == "" {
		zwsIDs = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	} else {
		for _, idStr := range strings.Split(zwsTypeIdsStr, ",") {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid zws_type_id"})
				return
			}
			zwsIDs = append(zwsIDs, id)
		}
	}

	categoryID := ctx.Param("category_id")
	categoryIDInt, err := strconv.Atoi(categoryID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный categoryID"})
		return
	}

	data, err := geodataRepository.GetPointsDataByZwsTypes(ctx, zwsIDs, categoryIDInt)

	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

// @Router /zulu/api/v2/points/categories/:category_id [get]
// @Summary Получение значений объектов
// @Tags Internal
// @Param zws_type_ids query []int false "ID типов объектов для фильтрации" collectionFormat(csv)
// @Param timestamp query string false "Timestamp для получения данных за конкретный период"
// @Produce  json
// @Success 200 {object} map[int]geodata.GetPointsDataByCategoryGroup
func getPointsCategoryByZwsTypeV2(ctx *gin.Context, geodataRepository geodata_repository.Repository) {
	zwsTypeIdsStr := ctx.Query("zws_type_id")

	var zwsIDs []int
	if zwsTypeIdsStr == "" {
		zwsIDs = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	} else {
		for _, idStr := range strings.Split(zwsTypeIdsStr, ",") {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid zws_type_id"})
				return
			}
			zwsIDs = append(zwsIDs, id)
		}
	}

	categoryID := ctx.Param("category_id")
	categoryIDInt, err := strconv.Atoi(categoryID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный categoryID"})
		return
	}

	timestamp := ctx.Query("timestamp")
	if len(timestamp) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный timestamp"})
		return
	}

	data, err := geodataRepository.GetPointsDataByZwsTypesV2(ctx, zwsIDs, categoryIDInt, timestamp)

	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}

	ctx.JSON(http.StatusOK, data)
}
