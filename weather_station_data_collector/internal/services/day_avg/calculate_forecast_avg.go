package day_avg

import (
	"context"
	"math"
	"time"
	"weather_station_data_collector/internal/models"
)

func (s service) CalculateForecastAvg(ctx context.Context) {
	rawForecast, err := s.ForecastApiClient.GetForecast(ctx)
	if err != nil {
		s.logInterrogateError(ctx, err)
		return
	}
	hourlyData := s.convertToWeatherData(*rawForecast)
	dailyData := make(map[string][]models.WeatherData)

	today := time.Now().Truncate(24 * time.Hour)

	for _, record := range hourlyData {
		date := record.Date.Format("2006-01-02")

		if record.Date.Truncate(24 * time.Hour).Equal(today) {
			continue
		}

		dailyData[date] = append(dailyData[date], record)
	}

	var dailyAverages []models.WeatherData
	for _, records := range dailyData {
		var totalTemp, totalHumidity, totalWindSpeed, totalWindGust, totalWindDirectionX, totalWindDirectionY, totalPressure float64
		var maxWindSpeed, minWindSpeed float64
		var count int

		minWindSpeed = records[0].WindSpeed
		maxWindSpeed = records[0].WindSpeed

		for _, record := range records {
			totalTemp += record.OutdoorTemperature
			totalHumidity += record.Humidity
			totalWindSpeed += record.WindSpeed
			totalWindGust += record.WindGustSpeed
			totalPressure += record.Pressure

			totalWindDirectionX += math.Cos(record.WindDirection * math.Pi / 180)
			totalWindDirectionY += math.Sin(record.WindDirection * math.Pi / 180)

			if record.WindSpeed > maxWindSpeed {
				maxWindSpeed = record.WindSpeed
			}
			if record.WindSpeed < minWindSpeed {
				minWindSpeed = record.WindSpeed
			}

			count++
		}

		avgWindDirectionX := totalWindDirectionX / float64(count)
		avgWindDirectionY := totalWindDirectionY / float64(count)

		avgWindDirection := math.Atan2(avgWindDirectionY, avgWindDirectionX) * 180 / math.Pi
		if avgWindDirection < 0 {
			avgWindDirection += 360
		}

		dailyAverages = append(dailyAverages, models.WeatherData{
			Date:               records[0].Date,
			OutdoorTemperature: totalTemp / float64(count),
			Humidity:           totalHumidity / float64(count),
			WindDirection:      avgWindDirection,
			WindSpeed:          totalWindSpeed / float64(count),
			MaxWindSpeed:       maxWindSpeed,
			MinWindSpeed:       minWindSpeed,
			WindGustSpeed:      totalWindGust / float64(count),
			Pressure:           totalPressure / float64(count),
		})
	}

	err = s.WeatherDataRepository.InsertWeatherDataButch(ctx, dailyAverages)
	if err != nil {
		s.logInterrogateError(ctx, err)
	} else {
		s.logInterrogateSuccess(ctx)
	}
}

func (s service) convertToWeatherData(originalData models.WeatherResponse) []models.WeatherData {
	var result []models.WeatherData
	for _, record := range originalData.List {
		dateTime := time.Unix(record.Dt, 0)
		data := models.WeatherData{
			Date: time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(),
				0, 0, 0, 0, time.UTC),
			OutdoorTemperature: record.Main.Temp,
			Humidity:           float64(record.Main.Humidity),
			WindDirection:      float64(record.Wind.Deg),
			WindSpeed:          record.Wind.Speed,
			MaxWindSpeed:       record.Wind.Speed,
			MinWindSpeed:       record.Wind.Speed,
			WindGustSpeed:      record.Wind.Gust,
			Pressure:           float64(record.Main.Pressure) * 0.75006375541921,
		}
		result = append(result, data)
	}
	return result
}
