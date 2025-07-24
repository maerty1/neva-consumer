package weather_data

import (
	"context"
	"fmt"
	"strings"
	"weather_station_data_collector/internal/models"
)

func (r *repository) InsertWeatherDataButch(ctx context.Context, data []models.WeatherData) error {
	query := `INSERT INTO "weather_data" (outdoor_temperature, humidity, wind_direction, wind_speed, max_wind_speed, 
                            min_wind_speed, wind_gust_speed, date, pressure)
VALUES `

	var values []string
	var args []interface{}
	for i, weather := range data {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			(i*9)+1, (i*9)+2, (i*9)+3, (i*9)+4, (i*9)+5, (i*9)+6, (i*9)+7, (i*9)+8, (i*9)+9))

		args = append(args, weather.OutdoorTemperature)
		args = append(args, weather.Humidity)
		args = append(args, weather.WindDirection)
		args = append(args, fmt.Sprintf("%d-%d {%d}", int(weather.MinWindSpeed), int(weather.MaxWindSpeed), int(weather.WindGustSpeed)))
		args = append(args, weather.MaxWindSpeed)
		args = append(args, weather.MinWindSpeed)
		args = append(args, weather.WindGustSpeed)
		args = append(args, weather.Date)
		args = append(args, weather.Pressure)
	}

	query += strings.Join(values, ", ")

	query += ` ON CONFLICT (date) 
DO UPDATE SET 
    outdoor_temperature = EXCLUDED.outdoor_temperature,
    humidity = EXCLUDED.humidity,
    wind_direction = EXCLUDED.wind_direction,
    wind_speed = EXCLUDED.wind_speed,
    max_wind_speed = EXCLUDED.max_wind_speed,
    min_wind_speed = EXCLUDED.min_wind_speed,
    wind_gust_speed = EXCLUDED.wind_gust_speed,
    pressure = EXCLUDED.pressure;`

	_, err := r.db.DB().Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
