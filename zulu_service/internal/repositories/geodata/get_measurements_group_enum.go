package geodata

import (
	"context"
	"database/sql"
	"log"
	"zulu_service/internal/config/errors"
	"zulu_service/internal/models/geodata"
)

func (r *repository) GetMeasurementGroupsEnum(ctx context.Context) (map[int]geodata.MeasurementGroupEnum, error) {
	measurementsGroups := make(map[int]geodata.MeasurementGroupEnum)

	query := `
SELECT
	id,
	group_front_desc,
	measurement_unit,
	mtc.front_description
FROM
	public.measurement_groups mg
	LEFT JOIN measurement_type_conversion mtc ON mg.in = mtc.measurement_type_id
	`

	rows, err := r.db.DB().Query(ctx, query)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return nil, errors.ErrInternalError
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name, unit string
		var convertedUnit sql.NullString
		if err := rows.Scan(&id, &name, &unit, &convertedUnit); err != nil {
			log.Printf("Ошибка сканирования строки: %v", err)
			return nil, errors.ErrInternalError
		}

		finalUnit := unit

		if convertedUnit.Valid {
			finalUnit = convertedUnit.String
		}

		measurementsGroups[id] = geodata.MeasurementGroupEnum{
			Name: name,
			Unit: finalUnit,
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при итерации по строкам: %v", err)
		return nil, errors.ErrInternalError
	}

	return measurementsGroups, nil
}
