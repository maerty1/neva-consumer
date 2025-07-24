package zulu

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"bff_service/internal/api_clients/core"
)

const nDays = 60
const parseDateFormat = "2006-01-02"

type FullPoint struct {
	Address string   `json:"address"`
	Title   string   `json:"title"`
	Packets []Packet `json:"packets"`
}

type Packet struct {
	Datetime               string        `json:"datetime"`
	Measurements           []Measurement `json:"measurements"`
	Iscopied               bool          `json:"iscopied"`
	IsDataCopied           bool          `json:"is_data_copied"`
	IsCalculatedDataCopied bool          `json:"is_calculated_data_copied"`
}

type Measurement struct {
	ID             string                    `json:"id"`
	Data           MeasurementData2          `json:"data"`
	CalculatedData MeasurementCalculatedData `json:"calculated_data"`
}

type MeasurementData2 struct {
	In  *float64 `json:"in"`
	Out *float64 `json:"out"`
}

type MeasurementCalculatedData struct {
	In  *float64 `json:"in"`
	Out *float64 `json:"out"`
}

// Temp
type TempMeasurement struct {
	Rn             int                       `json:"rn"`
	Data           MeasurementData2          `json:"data"`
	CalculatedData MeasurementCalculatedData `json:"calculated_data"`
}

// @Router /zulu/api/v2/points/{elem_id}/full [get]
// @Summary Получение данных для раскрытой карточки
// @Tags Points
// @Produce  json
// @Success 200 {object} FullPoint
func (f *facade) GetFullPoint(c *gin.Context) {
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

	type coeffStruct struct {
		lersCoeff *float64
		zuluCoeff *float64
	}
	coeffMap := make(map[string]coeffStruct)

	for _, packetMap := range zuluFullPoint.Packets {
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
			coeff := coeffMap[groupID]
			measurement := Measurement{
				ID: groupID,
				Data: MeasurementData2{
					In:  multiplyPointers(tempMeasurement.Data.In, coeff.lersCoeff),
					Out: multiplyPointers(tempMeasurement.Data.Out, coeff.lersCoeff),
					// In:  tempMeasurement.Data.In,
					// Out: tempMeasurement.Data.Out,
				},
				CalculatedData: MeasurementCalculatedData{
					In:  multiplyPointers(tempMeasurement.CalculatedData.In, coeff.zuluCoeff),
					Out: multiplyPointers(tempMeasurement.CalculatedData.Out, coeff.zuluCoeff),
					// In:  tempMeasurement.CalculatedData.In,
					// Out: tempMeasurement.CalculatedData.Out,
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
	fullPoint = fillMissingDateTime(fullPoint, true)

	c.JSON(http.StatusOK, fullPoint)
}

// fillMissingCalculatedData заполняет отсутствующие calculated_data из других пакетов
// enabled - флаг для включения/отключения логики
func fillMissingCalculatedData(fullPoint FullPoint, enabled bool) FullPoint {
	if !enabled {
		return fullPoint
	}

	// Сбор данных для заполнения
	calculatedDataMap := make(map[string]MeasurementCalculatedData)

	for _, packet := range fullPoint.Packets {
		for _, measurement := range packet.Measurements {
			if measurement.CalculatedData.In != nil || measurement.CalculatedData.Out != nil {
				calculatedDataMap[measurement.ID] = measurement.CalculatedData
			}
		}
	}

	// Заполнение отсутствующих calculated_data
	for i, packet := range fullPoint.Packets {
		for j, measurement := range packet.Measurements {
			if measurement.CalculatedData.In == nil && measurement.CalculatedData.Out == nil {
				if data, exists := calculatedDataMap[measurement.ID]; exists {
					fullPoint.Packets[i].Measurements[j].CalculatedData = data
				}
			}
		}
	}

	return fullPoint
}

// fillMissingDateTime заполняет отсутствующие Datetime из предшествующего пакета
// enabled - флаг для включения/отключения логики
func fillMissingDateTime(fullPoint FullPoint, enabled bool) FullPoint {
	if !enabled {
		return fullPoint
	}

	resultSlice := createResultSlice(fullPoint)
	startDate, _ := time.Parse(parseDateFormat, fullPoint.Packets[0].Datetime)
	resultSlice = append(resultSlice, fullPoint.Packets[0])

	for i := 1; i <= len(fullPoint.Packets)-1; i++ {
		endDate, _ := time.Parse(parseDateFormat, fullPoint.Packets[i].Datetime)
		datesDifference := (endDate.Sub(startDate).Hours() / 24) - 1
		if (int64(datesDifference)) >= 1 {
			sliceForming := formingDates(fullPoint.Packets[i-1].Measurements, int64(datesDifference), startDate)
			resultSlice = append(resultSlice, sliceForming...)
		}
		resultSlice = append(resultSlice, fullPoint.Packets[i])
		startDate = endDate
	}

	fullPoint.Packets = resultSlice
	return fullPoint
}

// Создаем временное хранилище для заполнения данными
func createResultSlice(fullPoint FullPoint) []Packet {
	startDate, _ := time.Parse(parseDateFormat, fullPoint.Packets[0].Datetime)
	endDate, _ := time.Parse(parseDateFormat, fullPoint.Packets[len(fullPoint.Packets)-1].Datetime)
	days := (endDate.Sub(startDate).Hours() / 24) + 1
	resultFullPoint := make([]Packet, 0, int64(days))
	return resultFullPoint
}

// Формирование отсутствующих дат и заполнение
func formingDates(measurement []Measurement, diff int64, startDate time.Time) []Packet {
	resultForming := make([]Packet, 0, diff)
	for diff > 0 {
		startDate = startDate.Add(24 * time.Hour)
		date := startDate.Format(parseDateFormat)
		addDate := Packet{Datetime: date, Measurements: measurement, Iscopied: true, IsDataCopied: true, IsCalculatedDataCopied: true}
		resultForming = append(resultForming, addDate)
		diff--
	}
	return resultForming
}
