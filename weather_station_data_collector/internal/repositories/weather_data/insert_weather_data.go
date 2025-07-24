package weather_data

import (
	"fmt"
	"weather_station_data_collector/internal/models"

	"golang.org/x/net/context"
)

func (r *repository) InsertWeatherData(ctx context.Context, data models.WeatherData) error {
	query := `INSERT INTO "weather_data" (outdoor_temperature, humidity, wind_direction, wind_speed, max_wind_speed,
		min_wind_speed, wind_gust_speed, date, pressure)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (date)
	DO UPDATE SET
	outdoor_temperature = EXCLUDED.outdoor_temperature,
		humidity = EXCLUDED.humidity,
		wind_direction = EXCLUDED.wind_direction,
		wind_speed = EXCLUDED.wind_speed,
		max_wind_speed = EXCLUDED.max_wind_speed,
		min_wind_speed = EXCLUDED.min_wind_speed,
		wind_gust_speed = EXCLUDED.wind_gust_speed,
		pressure = EXCLUDED.pressure;`

	_, err := r.db.DB().Exec(ctx, query,
		data.OutdoorTemperature,
		data.Humidity,
		data.WindDirection,
		fmt.Sprintf("%d-%d {%d}", int(data.MinWindSpeed), int(data.MaxWindSpeed), int(data.WindGustSpeed)),
		data.MaxWindSpeed,
		data.MinWindSpeed,
		data.WindGustSpeed,
		data.Date,
		data.Pressure,
	)
	if err != nil {
		return err
	}
	return nil
}
