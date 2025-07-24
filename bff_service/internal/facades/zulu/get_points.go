package zulu

import (
	"bff_service/internal/api_clients/core"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type GetPointsDataResponse struct {
	ElemID       int                                              `json:"elem_id"`
	Measurements map[string]core.GetPointsDataResponseMeasurement `json:"measurements"`
}

type MeasurementData struct {
	I *float64 `json:"i" example:"53.24"`
	O *float64 `json:"o" example:"45.12"`
}

type PointResponse struct {
	ElemID                 int                                           `json:"elem_id"`
	Title                  string                                        `json:"title" example:"Котельная 22"`
	Address                string                                        `json:"address" example:"Улица Пушкина 12"`
	MeasurementGroups      map[int]core.GetPointsDataResponseMeasurement `json:"measurement_groups"`
	Coordinates            []float64                                     `json:"coordinates" description:"[lat, lon]" example:"55.751244,37.618423"`
	HasAccident            bool                                          `json:"has_accident"`
	IsDataCopied           bool                                          `json:"is_data_copied"`
	IsCalculatedDataCopied bool                                          `json:"is_calculated_data_copied"`
	Type                   int                                           `json:"type"`
}

// @Router /zulu/api/v1/points [get]
// @Deprecated
// @Summary Получение значений объектов
// @Tags Points
// @Param zws_type_ids query []int false "ID типов объектов для фильтрации" collectionFormat(csv)
// @Param timestamp query string false "Timestamp для получения данных за конкретный период"
// @Produce  json
// @Success 200 {object} []PointResponse
func (f *facade) GetPoints(c *gin.Context) {
	zwsTypeIdsStr := c.Query("zws_type_id")

	var ids []int
	if zwsTypeIdsStr == "" {
		ids = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	} else {
		for _, idStr := range strings.Split(zwsTypeIdsStr, ",") {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid zws_type_id"})
				return
			}
			ids = append(ids, id)
		}
	}

	points, err := f.zuluApiClient.GetPoints(ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var response []PointResponse
	var batchRequests []core.GetPointsDataRequest

	for _, point := range points {
		measurementsRequest := make(map[string]core.GetPointsDataRequestMeasurement)
		for _, mg := range point.MeasurementGroups {
			measurementsRequest[mg.In] = core.GetPointsDataRequestMeasurement{
				I: mg.In,
				O: mg.Out,
			}
		}

		reqData := core.GetPointsDataRequest{
			ElemID:       point.ElemID,
			Measurements: measurementsRequest,
		}

		batchRequests = append(batchRequests, reqData)
	}

	// Выполнение пакетного запроса GetPointsData
	timestamp := c.Query("timestamp")
	dataResponses, err := f.coreDataApiClient.GetPointsData(batchRequests, timestamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка получения данных измерений: %v", err)})
		return
	}

	// Создание карты ответов для быстрого доступа по ElemID
	dataResponseMap := make(map[int]core.GetPointsDataResponse)
	for _, dataResp := range dataResponses {
		dataResponseMap[dataResp.ElemID] = dataResp
	}

	// Объединение данных измерений с базовыми данными поинтов
	for _, point := range points {
		dataResp, exists := dataResponseMap[point.ElemID]
		if !exists {
			fmt.Printf("Данные измерений отсутствуют для ElemID %d\n", point.ElemID)
		}

		// Обработка данных измерений
		measurementData := make(map[int]core.GetPointsDataResponseMeasurement)
		for groupID, mg := range point.MeasurementGroups {
			var measurement core.GetPointsDataResponseMeasurement

			if dataResp.ElemID == 3058 {
				fmt.Println(dataResp.IsDataCopied)
			}

			if m, ok := dataResp.Measurements[mg.In]; ok {
				measurement.I = m.I
				measurement.O = m.O
			}

			measurementData[groupID] = measurement
		}

		// Формирование конечного поинта
		pointResponse := PointResponse{
			ElemID:            point.ElemID,
			Title:             point.Title,
			Address:           point.Address,
			MeasurementGroups: measurementData,
			Coordinates:       point.Coordinates,
			HasAccident:       point.HasAccident,
			Type:              point.Type,
			IsDataCopied:      dataResp.IsDataCopied,
		}

		response = append(response, pointResponse)
	}

	c.JSON(http.StatusOK, response)
}
