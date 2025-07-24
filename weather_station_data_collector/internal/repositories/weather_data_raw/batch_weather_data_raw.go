package weather_data_raw

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"weather_station_data_collector/internal/models"
)

func (r *repository) InsertWeatherDataBatch(ctx context.Context, data []models.WeatherDataRaw) error {
	if len(data) == 0 {
		return errors.New("слайс пуст")
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`INSERT INTO "weather_data_raw" 
		("dateutc", "tempinf", "humidityin", "baromrelin", "baromabsin", "tempf", "humidity", "winddir", 
		"windspeedmph", "windgustmph", "maxdailygust", "solarradiation", "uv", "rainratein", 
		"eventrainin", "hourlyrainin", "dailyrainin", "weeklyrainin", "monthlyrainin", "yearlyrainin", 
		"totalrainin", "wh65batt", "freq") VALUES `)

	var params []interface{}
	paramIndex := 1

	for i, v := range data {
		if i > 0 {
			queryBuilder.WriteString(", ")
		}
		queryBuilder.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			paramIndex, paramIndex+1, paramIndex+2, paramIndex+3, paramIndex+4, paramIndex+5, paramIndex+6, paramIndex+7,
			paramIndex+8, paramIndex+9, paramIndex+10, paramIndex+11, paramIndex+12, paramIndex+13, paramIndex+14,
			paramIndex+15, paramIndex+16, paramIndex+17, paramIndex+18, paramIndex+19, paramIndex+20, paramIndex+21, paramIndex+22))
		paramIndex += 23

		params = append(params,
			v.Dateutc, v.Tempinf, v.Humidityin, v.Baromrelin, v.Baromabsin, v.Tempf, v.Humidity, v.Winddir,
			v.Windspeedmph, v.Windgustmph, v.Maxdailygust, v.Solarradiation, v.Uv, v.Rainratein, v.Eventrainin,
			v.Hourlyrainin, v.Dailyrainin, v.Weeklyrainin, v.Monthlyrainin, v.Yearlyrainin, v.Totalrainin,
			v.Wh65Batt, v.Freq)
	}

	query := queryBuilder.String()
	_, err := r.db.DB().Exec(ctx, query, params...)
	if err != nil {
		return err
	}
	return nil
}
