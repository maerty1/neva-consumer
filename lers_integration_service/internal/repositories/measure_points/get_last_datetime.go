package measure_points

import (
	"context"
	"database/sql"
	"fmt"
)

func (r *repository) GetLastMeasurePointDatetime(ctx context.Context, measurePointID int) (string, error) {
	var datetime sql.NullString
	// query := `
	// SELECT MIN(datetime) FROM (
	// 	SELECT datetime FROM measure_points_data
	// 	WHERE measure_point_id = $1
	// 	UNION ALL
	// 	SELECT datetime FROM measure_points_data_day
	// 	WHERE measure_point_id = $1
	// ) AS combined_datetimes`
	query := `
	SELECT MAX(datetime) - INTERVAL '30 days'
	FROM (
		SELECT datetime FROM measure_points_data_day
		WHERE measure_point_id = $1
	) AS combined_datetimes`
	err := r.db.DB().QueryRow(ctx, query, measurePointID).Scan(&datetime)
	if err != nil {
		if err.Error() != "no rows in result set" {
			return "", fmt.Errorf("не удалось получить самый старый datetime: %v", err)
		}
	}
	if !datetime.Valid {
		return "", nil
	}
	return datetime.String, nil
}
