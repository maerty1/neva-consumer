package day_avg

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
	"weather_station_data_collector/internal/models"
)

func (s service) CalculateCurrentDayAvg(ctx context.Context) {
	currentTime := time.Now()
	currentDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC)

	data, maxWindSpeed, minMinSpeed, err := s.RawWeatherDataRepository.SelectWeatherDataByDay(ctx, currentDay)
	if err != nil {
		s.logInterrogateError(ctx, errors.New(
			fmt.Sprintf(
				"Не удалось получить данные за день %s: %s",
				currentDay.Format("2006-01-02"), err.Error())))
	}
	weatherData := calculateAverages(data, currentDay)
	weatherData.MaxWindSpeed, weatherData.MinWindSpeed = maxWindSpeed, minMinSpeed
	weatherData.Pressure = weatherData.Pressure * 25.3285279248736
	weatherData.OutdoorTemperature = (weatherData.OutdoorTemperature - 32) * 1.8
	err = s.WeatherDataRepository.InsertWeatherData(ctx, weatherData)
	if err != nil {
		s.logInterrogateError(ctx, err)
	} else {
		s.logInterrogateSuccess(ctx)
	}
}

func calculateAverages(data []models.WeatherDataRaw, date time.Time) models.WeatherData {
	totals := make(map[string]float64)
	counts := make(map[string]int)

	for _, record := range data {
		processField(record.Tempf, "OutdoorTemperature", totals, counts)
		processField(record.Humidity, "Humidity", totals, counts)
		processField(record.Winddir, "WindDirection", totals, counts)
		processField(record.Windspeedmph, "WindSpeed", totals, counts)
		processField(record.Windgustmph, "WindGustSpeed", totals, counts)
		processField(record.Baromabsin, "Pressure", totals, counts)
	}

	return models.WeatherData{
		Date:               date,
		OutdoorTemperature: calculateAverage("OutdoorTemperature", totals, counts),
		Humidity:           calculateAverage("Humidity", totals, counts),
		WindDirection:      calculateAverage("WindDirection", totals, counts),
		WindSpeed:          calculateAverage("WindSpeed", totals, counts),
		WindGustSpeed:      calculateAverage("WindGustSpeed", totals, counts),
		Pressure:           calculateAverage("Pressure", totals, counts),
	}
}

func processField(value string, field string, totals map[string]float64, counts map[string]int) {
	if num, err := strconv.ParseFloat(value, 64); err == nil {
		totals[field] += num
		counts[field]++
	}
}

func calculateAverage(field string, totals map[string]float64, counts map[string]int) float64 {
	if counts[field] == 0 {
		return 0
	}
	return totals[field] / float64(counts[field])
}
