package zulu

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"

	"bff_service/internal/api_clients/core"
)

// @Router /zulu/api/v1/points/{elem_id}/full [get]
// @Summary Получение данных для раскрытой карточки
// @Tags Points
// @Produce  json
// @Success 200 {object} FullPoint
func (f *facade) GetPointsWithoutCopy(c *gin.Context) {
	elementID := c.Param("elem_id")
	elementIDInt, err := strconv.Atoi(elementID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный elementID"})
		return
	}

	zuluFullPoint, err := f.zuluApiClient.GetFullPoint(elementIDInt, nDays)
	if err != nil {
		// TODO: поменять на нормальные ошибки
		if err.Error() == "неожиданный статус-код: 404" {
			c.JSON(http.StatusNotFound, "Указанный elem_id не найден")
			return
		}
		log.Printf("Ошибка при получении данных из Zulu API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных из Zulu API"})
		return
	}
	measurementsToRequest := make(map[int]core.GetPointsDataHistoryMeasurementsRequest)

	for _, packetMap := range zuluFullPoint.Packets {
		for groupID, measurements := range packetMap {
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

	coreResponses, err := f.coreDataApiClient.GetPointsDataHistory(coreRequest, nDays, "")
	if err != nil {
		log.Printf("Ошибка при получении данных из Core API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных из Core API"})
		return
	}

	coreData, ok := coreResponses[strconv.Itoa(elementIDInt)]
	if !ok {
		fmt.Println("GOd damn")
	}

	tempData := make(map[string]*map[string]TempMeasurement)
	for ts, groups := range zuluFullPoint.Packets {
		data, ok := tempData[ts]

		if !ok {
			newData := make(map[string]TempMeasurement)
			data = &newData
			tempData[ts] = data
		}

		for groupID, groupData := range groups {
			(*data)[strconv.Itoa(groupID)] = TempMeasurement{
				CalculatedData: MeasurementCalculatedData{In: groupData.CalculatedData.In, Out: groupData.CalculatedData.Out},
			}
		}
	}

	for ts, groups := range coreData {
		data, ok := tempData[ts]

		if !ok {
			newData := make(map[string]TempMeasurement)
			data = &newData
			tempData[ts] = data
		}

		for groupID, groupData := range groups {
			tempMeasurement, exists := (*data)[groupID]
			if !exists {
				tempMeasurement = TempMeasurement{}
			}
			// Обновляем только поле Data, сохраняя CalculatedData
			tempMeasurement.Data = MeasurementData2{In: groupData.I, Out: groupData.O}
			(*data)[groupID] = tempMeasurement
		}
	}

	fullPoint := FullPoint{
		Address: zuluFullPoint.Address,
		Title:   zuluFullPoint.Title,
	}

	for packetTime, measurementsMap := range tempData {
		packet := Packet{
			Datetime:     packetTime,
			Measurements: []Measurement{},
		}
		for groupID, tempMeasurement := range *measurementsMap {

			measurement := Measurement{
				ID: groupID,
				Data: MeasurementData2{
					In:  tempMeasurement.Data.In,
					Out: tempMeasurement.Data.Out,
				},
				CalculatedData: MeasurementCalculatedData{
					In:  tempMeasurement.CalculatedData.In,
					Out: tempMeasurement.CalculatedData.Out,
				},
			}
			packet.Measurements = append(packet.Measurements, measurement)
		}

		fullPoint.Packets = append(fullPoint.Packets, packet)
	}

	sort.Slice(fullPoint.Packets, func(i, j int) bool {
		return fullPoint.Packets[i].Datetime < fullPoint.Packets[j].Datetime
	})

	fullPoint = fillMissingCalculatedData(fullPoint, true)

	c.JSON(http.StatusOK, fullPoint)
}
