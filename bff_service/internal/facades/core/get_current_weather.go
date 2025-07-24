package core

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WeatherData struct {
	PASSKEY        string `json:"PASSKEY"`
	StationType    string `json:"stationtype"`
	Runtime        string `json:"runtime"`
	DateUTC        string `json:"dateutc"`
	TempIndoorF    string `json:"tempinf"`
	HumidityIndoor string `json:"humidityin"`
	BaromRelIn     string `json:"baromrelin"`
	BaromAbsIn     string `json:"baromabsin"`
	TempOutdoorF   string `json:"tempf"`
	Humidity       string `json:"humidity"`
	WindDir        string `json:"winddir"`
	WindSpeedMPH   string `json:"windspeedmph"`
	WindGustMPH    string `json:"windgustmph"`
	MaxDailyGust   string `json:"maxdailygust"`
	SolarRadiation string `json:"solarradiation"`
	UV             string `json:"uv"`
	RainRateIn     string `json:"rainratein"`
	EventRainIn    string `json:"eventrainin"`
	HourlyRainIn   string `json:"hourlyrainin"`
	DailyRainIn    string `json:"dailyrainin"`
	WeeklyRainIn   string `json:"weeklyrainin"`
	MonthlyRainIn  string `json:"monthlyrainin"`
	YearlyRainIn   string `json:"yearlyrainin"`
	TotalRainIn    string `json:"totalrainin"`
	WH65Batt       string `json:"wh65batt"`
	Freq           string `json:"freq"`
	Model          string `json:"model"`
	Interval       string `json:"interval"`
}

// @Router /core/api/v1/weather/current [get]
// @Summary Получение текущей погоды
// @Description Получение текущей погоды
// @Tags Core
// @Produce  json
// @Success 200 {object} WeatherData "Успешный ответ"
func (f *facade) GetCurrentWeather(c *gin.Context) {
	url := "http://94.25.30.59:3000/weather/pull/last"
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get weather data"})
		return
	}
	defer resp.Body.Close()

	var weatherData WeatherData
	err = json.NewDecoder(resp.Body).Decode(&weatherData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to decode weather data"})
		return
	}

	c.JSON(http.StatusOK, weatherData)
}
