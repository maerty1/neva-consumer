package weather_data_raw

import (
	"context"
	"time"
	"weather_station_data_collector/internal/models"
)

func (r *repository) UpdateRawWeatherData(ctx context.Context, rawData models.WeatherDataRaw, oldTime time.Time) error {
	query := `
		UPDATE "weather_data_raw"
		SET 
			"dateutc" = $1,
			"tempinf" = CAST(NULLIF($2, '') AS double precision),
			"humidityin" = CAST(NULLIF($3, '') AS double precision),
			"baromrelin" = CAST(NULLIF($4, '') AS double precision),
			"baromabsin" = CAST(NULLIF($5, '') AS double precision),
			"tempf" = CAST(NULLIF($6, '') AS double precision),
			"humidity" = CAST(NULLIF($7, '') AS double precision),
			"winddir" = CAST(NULLIF($8, '') AS double precision),
			"windspeedmph" = CAST(NULLIF($9, '') AS double precision),
			"windgustmph" = CAST(NULLIF($10, '') AS double precision),
			"maxdailygust" = CAST(NULLIF($11, '') AS double precision),
			"solarradiation" = CAST(NULLIF($12, '') AS double precision),
			"uv" = CAST(NULLIF($13, '') AS double precision),
			"rainratein" = CAST(NULLIF($14, '') AS double precision),
			"eventrainin" = CAST(NULLIF($15, '') AS double precision),
			"hourlyrainin" = CAST(NULLIF($16, '') AS double precision),
			"dailyrainin" = CAST(NULLIF($17, '') AS double precision),
			"weeklyrainin" = CAST(NULLIF($18, '') AS double precision),
			"monthlyrainin" = CAST(NULLIF($19, '') AS double precision),
			"yearlyrainin" = CAST(NULLIF($20, '') AS double precision),
			"totalrainin" = CAST(NULLIF($21, '') AS double precision),
			"wh65batt" = CAST(NULLIF($22, '') AS double precision),
			"freq" = NULLIF($23, '')
		WHERE "dateutc" = $24;
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
		oldTime,
	)
	if err != nil {
		return err
	}
	return nil
}
