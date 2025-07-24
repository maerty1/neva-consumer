package measure_points

import (
	"context"
	"fmt"
	"strings"
)

type MeasurePointsDataDay struct {
	MeasurePointID int
	DateTime       string
	Values         string
}

func (r *repository) InsertMeasurePointDayDataBatch(ctx context.Context, data []MeasurePointsDataDay) error {
	if len(data) == 0 {
		return nil
	}

	query := `INSERT INTO measure_points_data_day (measure_point_id, datetime, values) VALUES `

	// Используем placeholders для каждого набора значений
	valueStrings := make([]string, 0, len(data))
	valueArgs := make([]interface{}, 0, len(data)*3)

	for i, d := range data {
		// PostgreSQL использует нумерацию placeholders начиная с 1
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		valueArgs = append(valueArgs, d.MeasurePointID, d.DateTime, d.Values)
	}

	query += strings.Join(valueStrings, ", ")
	query += ` ON CONFLICT (measure_point_id, datetime) DO UPDATE SET values = EXCLUDED.values`

	_, err := r.db.DB().Exec(ctx, query, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}
