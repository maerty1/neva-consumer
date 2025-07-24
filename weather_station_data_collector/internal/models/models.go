package models

import "time"

type WeatherStationResponseChannel struct {
	Code int
	Err  error
	Body WeatherDataRaw
}

type WeatherDataRaw struct {
	Dateutc        string `json:"dateutc"`
	Tempinf        string `json:"tempinf"`
	Humidityin     string `json:"humidityin"`
	Baromrelin     string `json:"baromrelin"`
	Baromabsin     string `json:"baromabsin"`
	Tempf          string `json:"tempf"`
	Humidity       string `json:"humidity"`
	Winddir        string `json:"winddir"`
	Windspeedmph   string `json:"windspeedmph"`
	Windgustmph    string `json:"windgustmph"`
	Maxdailygust   string `json:"maxdailygust"`
	Solarradiation string `json:"solarradiation"`
	Uv             string `json:"uv"`
	Rainratein     string `json:"rainratein"`
	Eventrainin    string `json:"eventrainin"`
	Hourlyrainin   string `json:"hourlyrainin"`
	Dailyrainin    string `json:"dailyrainin"`
	Weeklyrainin   string `json:"weeklyrainin"`
	Monthlyrainin  string `json:"monthlyrainin"`
	Yearlyrainin   string `json:"yearlyrainin"`
	Totalrainin    string `json:"totalrainin"`
	Wh65Batt       string `json:"wh65batt"`
	Freq           string `json:"freq"`
}

type WeatherData struct {
	Date               time.Time
	OutdoorTemperature float64
	Humidity           float64
	WindDirection      float64
	WindSpeed          float64
	MaxWindSpeed       float64
	MinWindSpeed       float64
	WindGustSpeed      float64
	Pressure           float64
}

type WeatherResponse struct {
	Cod     string        `json:"cod"`
	Message int           `json:"message"`
	Cnt     int           `json:"cnt"`
	List    []WeatherItem `json:"list"`
}

type WeatherItem struct {
	Dt      int64           `json:"dt"`
	Main    Main            `json:"main"`
	Weather []WeatherDetail `json:"weather"` // Это массив объектов
	Clouds  Clouds          `json:"clouds"`
	Wind    Wind            `json:"wind"`
	Rain    *Rain           `json:"rain,omitempty"` // Может отсутствовать
	Sys     Sys             `json:"sys"`
	DtTxt   string          `json:"dt_txt"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

type WeatherDetail struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Clouds struct {
	All int `json:"all"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

type Rain struct {
	ThreeH float64 `json:"3h"` // Используйте обратные кавычки
}

type Sys struct {
	Pod string `json:"pod"`
}
