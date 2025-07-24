package weather_data_raw

import (
	"context"
	"time"
)

func (r *repository) SelectTimeWithNullData(ctx context.Context) ([]time.Time, error) {
	query := `SELECT "dateutc"
FROM "weather_data_raw"
WHERE
    "tempinf" IS NULL AND "humidityin" IS NULL AND "baromrelin" IS NULL AND "baromabsin" IS NULL AND "tempf" IS NULL AND
    "humidity" IS NULL AND "winddir" IS NULL AND "windspeedmph" IS NULL AND
    "windgustmph" IS NULL AND "maxdailygust" IS NULL AND "solarradiation" IS NULL AND "uv" IS NULL AND "rainratein" IS NULL AND
    "eventrainin" IS NULL AND "hourlyrainin" IS NULL AND "dailyrainin" IS NULL AND "weeklyrainin"IS NULL AND
    "monthlyrainin" IS NULL AND "yearlyrainin"IS NULL AND "totalrainin"IS NULL AND "wh65batt"IS NULL AND "freq" IS NULL
`

	rows, err := r.db.DB().Query(ctx, query)
	var res []time.Time
	var temp time.Time
	for rows.Next() {
		err = rows.Scan(&temp)
		if err != nil {
			return nil, err
		}
		res = append(res, temp)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}
