package measure_points_data_day

import (
	"context"
	"time"
)

func (r repository) GetDataParameterByDay(ctx context.Context, day time.Time, paramName string) (float64, error) {
	query := `SELECT 
    (elem->>'value')::NUMERIC AS value
FROM 
    measure_points_data_day,
    LATERAL jsonb_array_elements(values) AS elem
WHERE 
    DATE(datetime) = $1
    AND elem->>'dataParameter' = $2
	AND measure_point_id = 779
ORDER BY datetime ASC;`
	row := r.db.DB().QueryRow(ctx, query, day, paramName)
	var res float64
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}
