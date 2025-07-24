package measure_points

import (
	"context"
)

func (r *repository) InsertMeasurePointData(ctx context.Context, measurePointID int, datetime string, values string) error {
	query := `
	INSERT INTO measure_points_data (measure_point_id, datetime, values) 
	VALUES ($1, $2, $3)
	ON CONFLICT (measure_point_id, datetime) DO UPDATE SET values = $3`
	_, err := r.db.DB().Exec(ctx, query, measurePointID, datetime, values)
	if err != nil {
		return err
	}

	return nil
}
