package models

import "time"

type WeatherData struct {
	Today struct {
		TempAvg      float64   `json:"temp_avg"`
		Temp         float64   `json:"temp"`
		PressureAvg  float64   `json:"pressure_avg"`
		Pressure     float64   `json:"pressure"`
		Humidity     float64   `json:"humidity"`
		HumidityAvg  float64   `json:"humidity_avg"`
		WindSpeed    float64   `json:"wind_speed"`
		WindSpeedAvg float64   `json:"wind_speed_avg"`
		Date         time.Time `json:"date"`
	} `json:"today"`
}
