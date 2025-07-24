package geodata

import (
	"context"
	"log"

	"zulu_service/internal/config/errors"
	"zulu_service/internal/models/geodata"
)

func (r *repository) GetElementDataByID(ctx context.Context, elementID int) ([]geodata.ElementData, error) {
	var elementData []geodata.ElementData

	query := `
	SELECT
		orc.inserted_ts,
		orc.val,
		orc.record_type,
		orc.parameter
	FROM
		zulu.object_records orc
	LEFT JOIN
		zulu.zulu_records_blacklist zrb
	ON
		orc.parameter = zrb.parameter
	WHERE
		orc.td IS NULL
		AND orc.elem_id = $1
		AND zrb.parameter IS NULL -- Исключаем записи, которые есть в черном списке
	ORDER BY
		orc.inserted_ts;
	`

	rows, err := r.db.DB().Query(ctx, query, elementID)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		return nil, errors.ErrInternalError
	}
	defer rows.Close()

	for rows.Next() {
		var data geodata.ElementData
		err := rows.Scan(&data.InsertedTS, &data.Val, &data.RecordType, &data.Parameter)
		if err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			return nil, errors.ErrInternalError
		}
		elementData = append(elementData, data)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Ошибка при итерации по строкам: %v", err)
		return nil, errors.ErrInternalError
	}

	if len(elementData) == 0 {
		return nil, errors.NotFoundWithDetails("ElementData", elementID)
	}

	return elementData, nil
}
