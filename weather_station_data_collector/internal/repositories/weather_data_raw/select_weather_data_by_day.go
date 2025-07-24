package weather_data_raw

import (
	"fmt"
	"time"
	"weather_station_data_collector/internal/models"

	"golang.org/x/net/context"
)

func (r *repository) SelectWeatherDataByDay(ctx context.Context, day time.Time) ([]models.WeatherDataRaw, float64, float64, error) {
	query := `
		SELECT dateutc, tempinf, humidityin, baromrelin, baromabsin, tempf, humidity, winddir, 
		       windspeedmph, windgustmph, maxdailygust, solarradiation, uv, rainratein, 
		       eventrainin, hourlyrainin, dailyrainin, weeklyrainin, monthlyrainin, yearlyrainin, 
		       totalrainin, wh65batt, freq,
		       MAX(windspeedmph::float) AS max_windspeed,
		       MIN(windspeedmph::float) AS min_windspeed
		FROM weather_data_raw
		WHERE dateutc >= $1::timestamp
		  AND dateutc < $1::timestamp + interval '1 day'
		  AND tempinf IS NOT NULL
		GROUP BY dateutc, tempinf, humidityin, baromrelin, baromabsin, tempf, humidity, winddir, 
		         windspeedmph, windgustmph, maxdailygust, solarradiation, uv, rainratein, 
		         eventrainin, hourlyrainin, dailyrainin, weeklyrainin, monthlyrainin, yearlyrainin, 
		         totalrainin, wh65batt, freq;`
	rows, err := r.db.DB().Query(ctx, query, day)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var res []models.WeatherDataRaw
	var maxWind, minWind float64
	firstRow := true

	for rows.Next() {
		var temp models.WeatherDataRaw
		var rowMaxWind, rowMinWind float64
		err = rows.Scan(
			&temp.Dateutc,
			&temp.Tempinf,
			&temp.Humidityin,
			&temp.Baromrelin,
			&temp.Baromabsin,
			&temp.Tempf,
			&temp.Humidity,
			&temp.Winddir,
			&temp.Windspeedmph,
			&temp.Windgustmph,
			&temp.Maxdailygust,
			&temp.Solarradiation,
			&temp.Uv,
			&temp.Rainratein,
			&temp.Eventrainin,
			&temp.Hourlyrainin,
			&temp.Dailyrainin,
			&temp.Weeklyrainin,
			&temp.Monthlyrainin,
			&temp.Yearlyrainin,
			&temp.Totalrainin,
			&temp.Wh65Batt,
			&temp.Freq,
			&rowMaxWind,
			&rowMinWind,
		)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("error scanning row: %w", err)
		}
		if firstRow {
			maxWind = rowMaxWind
			minWind = rowMinWind
			firstRow = false
		} else {
			if rowMaxWind > maxWind {
				maxWind = rowMaxWind
			}
			if rowMinWind < minWind {
				minWind = rowMinWind
			}
		}
		res = append(res, temp)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, 0, fmt.Errorf("row iteration error: %w", err)
	}
	return res, maxWind, minWind, nil
}
