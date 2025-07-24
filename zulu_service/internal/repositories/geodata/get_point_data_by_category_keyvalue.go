package geodata

import (
	"context"
	"database/sql"
	"log"
	"math"
	"strconv"
	"zulu_service/internal/config/errors"
	"zulu_service/internal/models/geodata"
)

func (r *repository) GetPointDataByCategoryKeyvalue(ctx context.Context, elemID int, categoryID int) (geodata.GetPointDataByCategoryKeyvalue, error) {
	pointData := geodata.GetPointDataByCategoryKeyvalue{}

	measurementGroupsQuery := `
	with recursive scan as (
			(		
	SELECT
		mt.id,
		mt.front_desc,
		mt.zulu_un,
		orc.val,
		mt.zulu_var,
		mt.scada_var,
		orc.inserted_ts,
		cmp.rn
	FROM 
		public.measurement_types mt
		JOIN public.category_to_measurement_type cmp 
		ON cmp.measurement_types_id = mt.id
		JOIN zulu.object_records orc 
		ON orc."parameter" = mt.zulu_var
	WHERE
		cmp.category_id = $1
		AND orc.td IS NULL
		AND orc.elem_id = $2
        AND fd <= now() at time zone 'utc'
    order by mt.id asc, (fd::date) desc, record_priority asc, fd desc
	limit 1	
		)
		union all (
			select r.*
			from scan 
			join lateral (
				SELECT
		mt.id,
		mt.front_desc,
		mt.zulu_un,
		orc.val,
		mt.zulu_var,
		mt.scada_var,
		orc.inserted_ts,
		cmp.rn
	FROM 
		public.measurement_types mt
		JOIN public.category_to_measurement_type cmp ON cmp.measurement_types_id = mt.id
		JOIN zulu.object_records orc ON orc."parameter" = mt.zulu_var
	WHERE
		cmp.category_id = $1
		AND orc.td IS NULL
		AND orc.elem_id = $2
        AND fd <= now() at time zone 'utc'
		and mt.id > scan.id
	order by mt.id asc, fd::date desc, record_priority asc, fd desc
	limit 1 
			) r 
			on true
		)
		)
		select 
			*
		from scan
		;
	`

	measurementGroupsRows, err := r.db.DB().Query(ctx, measurementGroupsQuery, categoryID, elemID)
	if err != nil {
		log.Printf("Ошибка выполнения запроса на получение групп измерений: %v", err)
		return pointData, errors.ErrInternalError
	}
	defer measurementGroupsRows.Close()

	for measurementGroupsRows.Next() {
		var measurement geodata.MeasurementKeyvalue
		var scadaVar string
		var zuluVar string
		var id, rn int
		var trash string
		var val sql.NullString

		err := measurementGroupsRows.Scan(
			&id,
			&measurement.Name,
			&measurement.Unit,
			&val,
			&zuluVar,
			&scadaVar,
			&trash,
			&rn,
		)

		if val.Valid {
			valStr := val.String
			if floatVal, err := strconv.ParseFloat(valStr, 64); err == nil {
				rounded := roundTo(floatVal, 3)
				if rounded == math.Trunc(rounded) {
					measurement.Value = strconv.FormatInt(int64(rounded), 10)
				} else {
					formattedValue := strconv.FormatFloat(rounded, 'f', 3, 64)
					measurement.Value = trimTrailingZeros(formattedValue)
				}
			} else {
				measurement.Value = valStr
			}

			if len(scadaVar) == 0 {
				measurement.Source = "zulu"
			} else {
				measurement.Source = "scada"
			}

			measurement.Rn = rn
			if err != nil {
				log.Printf("Ошибка сканирования строки групп измерений: %v", err)
				return pointData, errors.ErrInternalError
			}

			pointData.Measurements = append(pointData.Measurements, measurement)
		}

	}
	return pointData, nil

}

func roundTo(value float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Round(value*pow) / pow
}

// trimTrailingZeros удаляет лишние нули и точку, если необходимо
func trimTrailingZeros(s string) string {
	// Используем формат 'f' с точностью -1 для удаления лишних нулей
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s // Возвращаем исходную строку в случае ошибки
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}
