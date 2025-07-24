package zulu

import (
	"bff_service/internal/api_clients/core"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var typeOrder = map[int]int{
	1: 0,
	8: 1,
	3: 2,
}

// @Router /zulu/api/v1/points/with_rawdata [get]
// @Summary Получение значений объектов
// @Tags Points
// @Param zws_type_id query []int false "ID типов объектов для фильтрации" collectionFormat(csv)
// @Param timestamp query string false "Timestamp для получения данных за конкретный период"
// @Produce  json
// @Success 200 {object} []PointResponse
func (f *facade) GetPointsWithRawdata(c *gin.Context) {
	elementIDs, err := f.coreDataApiClient.GetElementIDs()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	zwsTypeIdsStr := c.Query("zws_type_id")

	var zwsIDs []int
	if zwsTypeIdsStr == "" || zwsTypeIdsStr == "0" {
		zwsIDs = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	} else {
		for _, idStr := range strings.Split(zwsTypeIdsStr, ",") {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid zws_type_id"})
				return
			}
			zwsIDs = append(zwsIDs, id)
		}
	}

	timestamp := c.Query("timestamp")
	points, err := f.zuluApiClient.GetFilteredPoints(elementIDs, zwsIDs, timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var response []PointResponse
	var batchRequests []core.GetPointsDataRequest
	coeffMap := make(map[string]*float64)

	for _, point := range points {
		measurementsRequest := make(map[string]core.GetPointsDataRequestMeasurement)
		for _, mg := range point.MeasurementGroups {
			coeffMap[mg.In] = mg.Coeff
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

	dataResponses, err := f.coreDataApiClient.GetPointsData(batchRequests, timestamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка получения данных измерений: %v", err)})
		return
	}

	dataResponseMap := make(map[int]core.GetPointsDataResponse)
	for _, dataResp := range dataResponses {
		dataResponseMap[dataResp.ElemID] = dataResp
	}

	for _, point := range points {
		dataResp, exists := dataResponseMap[point.ElemID]
		if !exists {
			fmt.Printf("Данные измерений отсутствуют для ElemID %d\n", point.ElemID)
		}

		measurementData := make(map[int]core.GetPointsDataResponseMeasurement)
		for groupID, mg := range point.MeasurementGroups {
			var measurement core.GetPointsDataResponseMeasurement

			if m, ok := dataResp.Measurements[mg.In]; ok {
				coeff := coeffMap[mg.In]

				measurement.I = multiplyPointers(m.I, coeff)
				measurement.O = multiplyPointers(m.O, coeff)
			}

			measurementData[groupID] = measurement
		}

		pointResponse := PointResponse{
			ElemID:                 point.ElemID,
			Title:                  point.Title,
			Address:                point.Address,
			MeasurementGroups:      measurementData,
			Coordinates:            point.Coordinates,
			HasAccident:            point.HasAccident,
			Type:                   point.Type,
			IsDataCopied:           dataResp.IsDataCopied,
			IsCalculatedDataCopied: point.IsCopied,
		}

		response = append(response, pointResponse)
	}

	// sort.Slice(response, func(i, j int) bool {
	// 	orderI, foundI := typeOrder[response[i].Type]
	// 	orderJ, foundJ := typeOrder[response[j].Type]
	// 	if foundI && foundJ {
	// 		return orderI < orderJ
	// 	}

	// 	if foundI {
	// 		return true
	// 	}
	// 	if foundJ {
	// 		return false
	// 	}

	// 	return response[i].Type < response[j].Type
	// })

	sort.Slice(response, func(i, j int) bool {
		orderI, foundI := typeOrder[response[i].Type]
		orderJ, foundJ := typeOrder[response[j].Type]

		// Сначала сортируем по typeOrder, если оба Type присутствуют в typeOrder
		if foundI && foundJ {
			if orderI != orderJ {
				return orderI < orderJ
			}
			// Если Type одинаковы, сортируем по ElemID
			return response[i].ElemID < response[j].ElemID
		}

		// Если только один из Type присутствует в typeOrder, тот что присутствует идет первым
		if foundI {
			return true
		}
		if foundJ {
			return false
		}

		// Если ни один из Type не присутствует в typeOrder, сортируем по Type, а затем по ElemID
		if response[i].Type != response[j].Type {
			return response[i].Type < response[j].Type
		}
		return response[i].ElemID < response[j].ElemID
	})

	c.JSON(http.StatusOK, response)
}

func multiplyPointers(a, b *float64) *float64 {
	if a == nil {
		return nil
	} else if b == nil {
		return a
	}
	result := *a * *b
	result = roundToOneDecimalPlace(result)
	fmt.Println(result)
	return &result
}

func roundToOneDecimalPlace(value float64) float64 {
	return math.Round(value*10) / 10
}
