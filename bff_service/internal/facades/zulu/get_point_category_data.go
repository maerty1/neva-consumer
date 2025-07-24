package zulu

import (
	"bff_service/internal/api_clients/core"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Данные для получения сырья
type GroupMeasurementsData struct {
	In  *float64 `json:"in"`
	Out *float64 `json:"out"`
}

// Данные из Зулу
type GroupMeasurementsCalculatedData struct {
	In  *float64 `json:"in"`
	Out *float64 `json:"out"`
}

type GroupMeasurement struct {
	ID             string                          `json:"id"`
	CalculatedData GroupMeasurementsCalculatedData `json:"calculated_data"`
	Data           GroupMeasurementsData           `json:"data"`
	Rn             int                             `json:"rn"`
}

type GetPointDataByCategoryGroup struct {
	Measurements           []GroupMeasurement `json:"measurements"`
	IsDataCopied           bool               `json:"is_data_copied"`
	IsCalculatedDataCopied bool               `json:"is_calculated_data_copied"`
}

// @Router /zulu/api/v1/points/{elem_id}/categories/{category_id} [get]
// @Summary Получение данных для категорий
// @Param type query string true "Type of data" Enums(group, keyvalue)
// @Param timestamp query string false "Timestamp для получения данных за конкретный период"
// @Tags Points
// @Produce  json
// @Success 200 {object} GetPointDataByCategoryGroup "Ответ для type=group"
// @Success 201 {object} zulu.GetPointDataByCategoryKeyvalue "Ответ для type=keyvalue"
func (f *facade) GetPointCategoryData(c *gin.Context) {
	elementID := c.Param("elem_id")
	elementIDInt, err := strconv.Atoi(elementID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный elementID"})
		return
	}

	categoryID := c.Param("category_id")
	categoryIDInt, err := strconv.Atoi(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный categoryID"})
		return
	}
	groupType := c.Query("type")
	timestamp := c.Query("timestamp")

	var requestedTime time.Time
	if timestamp != "" {

		requestedTime, err = time.Parse(time.RFC3339, timestamp)
		if err != nil {
			layout := "2006-01-02"
			requestedTime, err = time.Parse(layout, timestamp)
			if err != nil {
				log.Printf("Неверный формат timestamp: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат timestamp"})
				return
			}
		}
	} else {
		requestedTime = time.Now()
	}

	if groupType == "group" {
		zuluFullPoint, err := f.zuluApiClient.GetPointCategoryDataGroup(elementIDInt, categoryIDInt, timestamp)
		if err != nil {
			log.Printf("Ошибка при получении данных из Zulu API: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных из Zulu API"})
			return
		}

		measurementsToRequest := make(map[int]core.GetPointsDataHistoryMeasurementsRequest)
		type coeffStruct struct {
			lersCoeff *float64
			zuluCoeff *float64
		}
		coeffMap := make(map[string]coeffStruct)

		for _, packetMap := range zuluFullPoint.Measurements {
			for groupID, measurements := range packetMap {
				groupIDstring := strconv.Itoa(groupID)

				coeffMap[groupIDstring] = coeffStruct{
					lersCoeff: measurements.LersCoeff,
					zuluCoeff: measurements.ZuluCoeff,
				}
				if measurements.Data.In != "" || measurements.Data.Out != "" {
					measurementsToRequest[groupID] = core.GetPointsDataHistoryMeasurementsRequest{
						I: measurements.Data.In,
						O: measurements.Data.Out,
					}
				}

			}
		}

		coreRequest := []core.GetPointsDataHistoryRequest{}
		coreRequest = append(coreRequest, core.GetPointsDataHistoryRequest{
			ElemID:       elementIDInt,
			Measurements: measurementsToRequest,
		})

		coreResponses, err := f.coreDataApiClient.GetPointsDataHistory(coreRequest, 1, timestamp)
		if err != nil {
			log.Printf("Ошибка при получении данных из Core API: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных из Core API"})
			return
		}

		ourData, ok := coreResponses[strconv.Itoa(elementIDInt)]
		if !ok {
			fmt.Println("GOd damn")
		}

		tempData := make(map[string]*map[string]TempMeasurement)
		var isCalculatedDataCopied bool

		for ts, groups := range zuluFullPoint.Measurements {
			data, ok := tempData[ts]

			if copied, err := IsDataCopied(ts, requestedTime); err != nil {
				log.Printf("Ошибка при сравнении временных меток: %v", err)
				copied = true
			} else if copied {
				isCalculatedDataCopied = true
			}

			if !ok {
				newData := make(map[string]TempMeasurement)
				data = &newData
				tempData[ts] = data
			}

			for groupID, groupData := range groups {
				(*data)[strconv.Itoa(groupID)] = TempMeasurement{
					CalculatedData: MeasurementCalculatedData{In: groupData.CalculatedData.In, Out: groupData.CalculatedData.Out}, Rn: groupData.Rn,
				}
			}
		}

		var isDataCopied bool
		for ts, groups := range ourData {
			data, ok := tempData[ts]

			if !ok {
				newData := make(map[string]TempMeasurement)
				data = &newData
				tempData[ts] = data
			}

			if copied, err := IsDataCopied(ts, requestedTime); err != nil {
				log.Printf("Ошибка при сравнении временных меток: %v", err)
				copied = true
			} else if copied {
				isDataCopied = true
			}

			for groupID, groupData := range groups {
				(*data)[groupID] = TempMeasurement{
					Data: MeasurementData2{In: groupData.I, Out: groupData.O},
				}
			}
		}

		fullPoint := GetPointDataByCategoryGroup{
			Measurements: make([]GroupMeasurement, 0),
		}

		for _, measurementsMap := range tempData {
			for groupID, tempMeasurement := range *measurementsMap {
				coeff := coeffMap[groupID]

				measurement := GroupMeasurement{
					ID: groupID,
					Rn: tempMeasurement.Rn,
					CalculatedData: GroupMeasurementsCalculatedData{
						In:  multiplyPointers(tempMeasurement.CalculatedData.In, coeff.zuluCoeff),
						Out: multiplyPointers(tempMeasurement.CalculatedData.Out, coeff.zuluCoeff),
					},
				}
				fullPoint.Measurements = append(fullPoint.Measurements, measurement)
			}
		}

		// Используем карту для объединения измерений по id
		mergedMeasurements := make(map[string]GroupMeasurement)
		for _, measurementsMap := range tempData {
			for groupID, tempMeasurement := range *measurementsMap {
				coeff := coeffMap[groupID]

				if existing, exists := mergedMeasurements[groupID]; exists {
					if tempMeasurement.CalculatedData.In != nil {
						existing.CalculatedData.In = tempMeasurement.CalculatedData.In
					}
					if tempMeasurement.CalculatedData.Out != nil {
						existing.CalculatedData.Out = tempMeasurement.CalculatedData.Out
					}
					if tempMeasurement.Data.In != nil {
						existing.Data.In = multiplyPointers(tempMeasurement.Data.In, coeff.lersCoeff)
					}
					if tempMeasurement.Data.Out != nil {
						existing.Data.Out = multiplyPointers(tempMeasurement.Data.Out, coeff.lersCoeff)

					}
					if tempMeasurement.Rn != 0 {
						existing.Rn = tempMeasurement.Rn
					}
					mergedMeasurements[groupID] = existing
				} else {
					mergedMeasurements[groupID] = GroupMeasurement{
						ID: groupID,
						Rn: tempMeasurement.Rn,
						Data: GroupMeasurementsData{
							In:  multiplyPointers(tempMeasurement.Data.In, coeff.lersCoeff),
							Out: multiplyPointers(tempMeasurement.Data.Out, coeff.lersCoeff),
						},
						CalculatedData: GroupMeasurementsCalculatedData{
							In:  tempMeasurement.CalculatedData.In,
							Out: tempMeasurement.CalculatedData.Out,
						},
					}
				}
			}
		}

		fullPoint = GetPointDataByCategoryGroup{
			Measurements: make([]GroupMeasurement, 0, len(mergedMeasurements)), IsDataCopied: isDataCopied, IsCalculatedDataCopied: isCalculatedDataCopied,
		}

		for _, measurement := range mergedMeasurements {
			fullPoint.Measurements = append(fullPoint.Measurements, measurement)
		}

		sort.Slice(fullPoint.Measurements, func(i, j int) bool {
			return fullPoint.Measurements[i].Rn < fullPoint.Measurements[j].Rn
		})
		c.JSON(http.StatusOK, fullPoint)
		return

	} else if groupType == "keyvalue" {
		pointData, err := f.zuluApiClient.GetPointCategoryDataKeyvalue(elementIDInt, categoryIDInt)
		if err != nil {
			log.Printf("Ошибка при получении данных из Zulu API: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных из Zulu API"})
			return
		}
		sort.Slice(pointData.Measurements, func(i, j int) bool {
			return pointData.Measurements[i].Rn < pointData.Measurements[j].Rn
		})
		c.JSON(http.StatusOK, pointData)
		return
	} else {
		c.JSON(http.StatusBadRequest, "неверный query param 'type'. Доступно: keyvalue, group")
	}
}

func IsDataCopied(dataTimestamp string, requestedTime time.Time) (bool, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05Z07:00",
	}

	var dataTime time.Time
	var err error

	for _, layout := range layouts {
		dataTime, err = time.Parse(layout, dataTimestamp)
		if err == nil {
			break
		}
	}
	if err != nil {
		return false, fmt.Errorf("invalid dataTimestamp format: %v", err)
	}

	return dataTime.Before(requestedTime), nil
}
