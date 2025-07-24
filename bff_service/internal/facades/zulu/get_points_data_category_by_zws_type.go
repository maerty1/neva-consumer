package zulu

import (
	"bff_service/internal/api_clients/core"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type CalculatedData struct {
	In  *float64 `json:"in"`
	Out *float64 `json:"out"`
}

type Data struct {
	In  *float64 `json:"in"`
	Out *float64 `json:"out"`
}

type GetPointsDataCategoryMeasurement struct {
	CalculatedData CalculatedData `json:"calculated_data"`
	Data           Data           `json:"data"`
}

type CategoryMeasurements struct {
	IsDataCopied           bool                                        `json:"is_data_copied"`
	IsCaclulatedDataCopied bool                                        `json:"is_calculated_data_copied"`
	Measurements           map[string]GetPointsDataCategoryMeasurement `json:"measurements"`
}

// Тип для всего ответа, где ключ — ID категории
type GetPointsDataCategoryResponse map[string]CategoryMeasurements

// @Router /zulu/api/v2/points/categories/{category_id} [get]
// @Summary Получение значений объектов
// @Tags Points
// @Param zws_type_ids query []int false "ID типов объектов для фильтрации" collectionFormat(csv)
// @Param timestamp query string false "Timestamp для получения данных за конкретный период"
// @Produce  json
// @Success 200 {object} GetPointsDataCategoryResponse
func (f *facade) GetPointsDataCategoryByZwsType(c *gin.Context) {
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

	categoryID := c.Param("category_id")
	categoryIDInt, err := strconv.Atoi(categoryID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный categoryID"})
		return
	}

	timestamp := c.Query("timestamp")
	if len(timestamp) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный timestamp"})
		return
	}

	points, err := f.zuluApiClient.GetPointsDataCategoryByZwsType(ids, categoryIDInt, timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var response GetPointsDataCategoryResponse = make(GetPointsDataCategoryResponse)
	var batchRequests []core.GetPointsDataRequest

	type coeffStruct struct {
		lersCoeff *float64
		zuluCoeff *float64
	}
	coeffMap := make(map[string]coeffStruct)

	for elemID, point := range points {
		measurementsRequest := make(map[string]core.GetPointsDataRequestMeasurement)
		for groupID, mg := range point.Measurements {
			coeffMap[groupID] = coeffStruct{
				lersCoeff: mg.LersCoeff,
				zuluCoeff: mg.ZuluCoeff,
			}

			measurementsRequest[groupID] = core.GetPointsDataRequestMeasurement{
				I: *mg.Data.In,
				O: *mg.Data.Out,
			}
		}

		elemIDint, _ := strconv.Atoi(elemID)

		reqData := core.GetPointsDataRequest{
			ElemID:       elemIDint,
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

	for elemID, point := range points {
		elemIDint, _ := strconv.Atoi(elemID)
		dataResp, exists := dataResponseMap[elemIDint]
		if !exists {
			fmt.Printf("Данные измерений отсутствуют для ElemID %v\n", elemID)
		}

		cm := CategoryMeasurements{
			Measurements: make(map[string]GetPointsDataCategoryMeasurement),
		}

		for groupID := range point.Measurements {
			group := GetPointsDataCategoryMeasurement{}
			if m, ok := dataResp.Measurements[groupID]; ok {
				coeff := coeffMap[groupID]

				group.Data = Data{
					In:  multiplyPointers(m.I, coeff.zuluCoeff),
					Out: multiplyPointers(m.O, coeff.zuluCoeff),
				}
			}

			calculatedData := CalculatedData(points[elemID].Measurements[groupID].CalculatedData)
			group.CalculatedData = calculatedData

			cm.Measurements[groupID] = group
			cm.IsDataCopied = dataResp.IsDataCopied
			cm.IsCaclulatedDataCopied = points[elemID].IsCalculatedCopied
		}

		response[elemID] = cm

	}

	c.JSON(http.StatusOK, response)
}
