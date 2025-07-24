package weather_data_raw

import (
	"context"
	"weather_station_data_collector/internal/models"
)

func (r *repository) InsertRawWeatherData(ctx context.Context, rawData models.WeatherDataRaw) error {
	query := `INSERT INTO "weather_data_raw" 
("dateutc", "tempinf", "humidityin", "baromrelin", "baromabsin", "tempf", "humidity", "winddir", "windspeedmph", 
 "windgustmph", "maxdailygust", "solarradiation", "uv", "rainratein", "eventrainin", "hourlyrainin", 
 "dailyrainin", "weeklyrainin", "monthlyrainin", "yearlyrainin", "totalrainin", "wh65batt", "freq")
VALUES ($1, CAST(NULLIF($2, '') AS double precision), CAST(NULLIF($3, '') AS double precision), 
        CAST(NULLIF($4, '') AS double precision), CAST(NULLIF($5, '') AS double precision), 
        CAST(NULLIF($6, '') AS double precision), CAST(NULLIF($7, '') AS double precision), 
        CAST(NULLIF($8, '') AS double precision), CAST(NULLIF($9, '') AS double precision), 
        CAST(NULLIF($10, '') AS double precision), CAST(NULLIF($11, '') AS double precision), 
        CAST(NULLIF($12, '') AS double precision), CAST(NULLIF($13, '') AS double precision), 
        CAST(NULLIF($14, '') AS double precision), CAST(NULLIF($15, '') AS double precision), 
        CAST(NULLIF($16, '') AS double precision), CAST(NULLIF($17, '') AS double precision), 
        CAST(NULLIF($18, '') AS double precision), CAST(NULLIF($19, '') AS double precision), 
        CAST(NULLIF($20, '') AS double precision), CAST(NULLIF($21, '') AS double precision), 
        CAST(NULLIF($22, '') AS double precision), NULLIF($23, ''));

`

	_, err := r.db.DB().Exec(ctx,
		query,
		rawData.Dateutc,
		rawData.Tempinf,
		rawData.Humidityin,
		rawData.Baromrelin,
		rawData.Baromabsin,
		rawData.Tempf,
		rawData.Humidity,
		rawData.Winddir,
		rawData.Windspeedmph,
		rawData.Windgustmph,
		rawData.Maxdailygust,
		rawData.Solarradiation,
		rawData.Uv,
		rawData.Rainratein,
		rawData.Eventrainin,
		rawData.Hourlyrainin,
		rawData.Dailyrainin,
		rawData.Weeklyrainin,
		rawData.Monthlyrainin,
		rawData.Yearlyrainin,
		rawData.Totalrainin,
		rawData.Wh65Batt,
		rawData.Freq,
	)
	if err != nil {
		return err
	}
	return nil
}
