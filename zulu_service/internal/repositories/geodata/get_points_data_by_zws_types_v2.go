package geodata

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"
	"zulu_service/internal/config/errors"
	"zulu_service/internal/models/geodata"
)

func (r *repository) GetPointsDataByZwsTypesV2(ctx context.Context, zwsTypeIDs []int, categoryID int, timestamp string) (map[int]*geodata.GetPointsDataByCategoryGroup, error) {
	var elemIDs []int

	var requestedTime time.Time
	var err error

	requestedTime, err = time.Parse(time.RFC3339, timestamp)
	if err != nil {
		layout := "2006-01-02"
		requestedTime, err = time.Parse(layout, timestamp)
		if err != nil {
			log.Printf("Неверный формат timestamp: %v", err)
			return nil, errors.ErrInvalidInput
		}
	}

	elemIdsQuery := `
    SELECT
        elem_id
    FROM zulu.elems_metadata
    WHERE zws_type_id = ANY($1)
    `

	rows, err := r.db.DB().Query(ctx, elemIdsQuery, zwsTypeIDs)
	if err != nil {
		log.Printf("ошибка выполнения elemIdsQuery: %v", err)
		return nil, errors.ErrInternalError
	}
	defer rows.Close()

	for rows.Next() {
		var elemID int
		err := rows.Scan(&elemID)
		if err != nil {
			log.Printf("ошибка сканирования elem id: %v", err)
			return nil, errors.ErrInternalError
		}
		elemIDs = append(elemIDs, elemID)
	}

	if err = rows.Err(); err != nil {
		log.Printf("ошибка итерации elem id rows: %v", err)
		return nil, errors.ErrInternalError
	}

	fullDataMap := make(map[int]*geodata.GetPointsDataByCategoryGroup)

	for _, elemID := range elemIDs {
		fullDataMap[elemID] = &geodata.GetPointsDataByCategoryGroup{
			Measurements: make(map[int]*geodata.GroupMeasurement),
			IsCopied:     false, // Инициализация флага
		}
	}

	measurementGroupsQuery := `
SELECT
    mg.id,
    ctg.rn,
    mg.group_front_desc,
    mg.measurement_unit,
    mg.in,
    mg.out,
    in_mt.rest_var AS in_lers_var,
    out_mt.rest_var AS out_lers_var,
    mtc.lers_coeff,
    mtc.zulu_coeff
FROM
    public.category_to_group ctg
    JOIN public.measurement_groups mg ON mg.id = ctg.group_id
    LEFT JOIN public.measurement_types in_mt ON mg.in = in_mt.id
    LEFT JOIN public.measurement_types out_mt ON mg.out = out_mt.id
    LEFT JOIN public.measurement_type_conversion mtc ON mtc.measurement_type_id = in_mt.id
WHERE
    ctg.category_id = $1
    `

	rows, err = r.db.DB().Query(ctx, measurementGroupsQuery, categoryID)
	if err != nil {
		log.Printf("ошибка выполнения measurement groups query: %v", err)
		return nil, errors.ErrInternalError
	}
	defer rows.Close()

	type measurementGroup struct {
		ID                   int
		Rn                   int
		Name                 string
		Unit                 string
		InMeasurementTypeID  int
		OutMeasurementTypeID int
		InLersVariable       string
		OutLersVariable      string
		LersCoeff            sql.NullFloat64
		ZuluCoeff            sql.NullFloat64
	}

	var measurementGroups []measurementGroup
	measurementGroupMap := make(map[int]measurementGroup)

	for rows.Next() {
		var mg measurementGroup
		err := rows.Scan(
			&mg.ID,
			&mg.Rn,
			&mg.Name,
			&mg.Unit,
			&mg.InMeasurementTypeID,
			&mg.OutMeasurementTypeID,
			&mg.InLersVariable,
			&mg.OutLersVariable,
			&mg.LersCoeff,
			&mg.ZuluCoeff,
		)
		if err != nil {
			log.Printf("ошибка сканирования measurement group row: %v", err)
			return nil, errors.ErrInternalError
		}
		measurementGroups = append(measurementGroups, mg)
		measurementGroupMap[mg.ID] = mg
	}
	if err = rows.Err(); err != nil {
		log.Printf("ошибка итерации measurement group rows: %v", err)
		return nil, errors.ErrInternalError
	}

	if len(measurementGroups) == 0 {
		return fullDataMap, nil
	}

	relevantGroupIDs := make([]int, 0, len(measurementGroups))
	for _, mg := range measurementGroups {
		relevantGroupIDs = append(relevantGroupIDs, mg.ID)
	}

	measureQuery := `
WITH all_combinations AS (
    SELECT 
        e.elem_id,
        mg.id AS measurement_group_id,
        mg.group_front_desc AS measurement_group_name,
        mg.measurement_unit,
        mg.in AS in_measurement_type_id,
        mg.out AS out_measurement_type_id,
        in_mt.rest_var AS in_lers_var,
        out_mt.rest_var AS out_lers_var
    FROM 
        unnest($1::int[]) AS e(elem_id)
    CROSS JOIN 
        public.measurement_groups mg
    LEFT JOIN 
        public.measurement_types in_mt ON mg.in = in_mt.id
    LEFT JOIN 
        public.measurement_types out_mt ON mg.out = out_mt.id
    WHERE 
        mg.id = ANY($2)
),
latest_records AS (
    SELECT 
        ac.elem_id,
        ac.measurement_group_id,
        ac.measurement_group_name,
        ac.measurement_unit,
        mt.id AS measurement_type_id,
        mt.front_desc AS measurement_type_name,
        mt.zulu_un AS measurement_type_unit,
        CASE
            WHEN mt.id = ac.in_measurement_type_id THEN 'in'
            WHEN mt.id = ac.out_measurement_type_id THEN 'out'
            ELSE 'unknown'
        END AS measurement_role,
        orc.val,
        mt.rest_var,
        mt.zulu_var,
        orc.inserted_ts, 
        ROW_NUMBER() OVER (
            PARTITION BY 
                ac.elem_id, 
                ac.measurement_group_id, 
                CASE
                    WHEN mt.id = ac.in_measurement_type_id THEN 'in'
                    WHEN mt.id = ac.out_measurement_type_id THEN 'out'
                    ELSE 'unknown'
                END
            ORDER BY 
                orc.inserted_ts DESC
        ) AS rn
    FROM 
        all_combinations ac
    LEFT JOIN 
        zulu.object_records orc 
        ON orc.elem_id = ac.elem_id
        AND orc.inserted_ts <= $3::timestamp + INTERVAL '1 day'
        AND orc.td IS NULL
    LEFT JOIN 
        public.measurement_types mt 
        ON mt.id = ac.in_measurement_type_id OR mt.id = ac.out_measurement_type_id
)
SELECT
    lr.elem_id,
    lr.measurement_group_id,
    lr.measurement_group_name,
    lr.measurement_unit,
    lr.measurement_type_id,
    lr.measurement_type_name,
    lr.measurement_type_unit,
    lr.measurement_role,
    lr.val,
    lr.rest_var,
    lr.zulu_var,
    lr.inserted_ts -- Добавляем inserted_ts в итоговый SELECT
FROM
    latest_records lr
WHERE
    lr.rn = 1;
    `

	rows, err = r.db.DB().Query(ctx, measureQuery, elemIDs, relevantGroupIDs, requestedTime)
	if err != nil {
		log.Printf("ошибка выполнения measurements query: %v", err)
		return nil, errors.ErrInternalError
	}
	defer rows.Close()

	for rows.Next() {
		var (
			elemID          int
			mgID            int
			mgName          string
			mgUnit          string
			mtID            int
			mtName          sql.NullString
			mtUnit          sql.NullString
			measurementRole string
			val             sql.NullString
			lersVariable    sql.NullString
			zuluVariable    sql.NullString
			insertedTS      sql.NullTime
		)

		err := rows.Scan(
			&elemID,
			&mgID,
			&mgName,
			&mgUnit,
			&mtID,
			&mtName,
			&mtUnit,
			&measurementRole,
			&val,
			&lersVariable,
			&zuluVariable,
			&insertedTS, // Сканируем поле inserted_ts
		)
		if err != nil {
			log.Printf("ошибка сканирования measurement row: %v", err)
			continue
		}

		var parsedVal *float64
		if val.Valid {
			f, err := strconv.ParseFloat(val.String, 64)
			if err != nil {
				log.Printf("ошибка парсинга val '%s' для elem_id %d, mg_id %d: %v", val.String, elemID, mgID, err)
				parsedVal = nil
			} else {
				parsedVal = &f
			}
		} else {
			parsedVal = nil
		}

		fullData, exists := fullDataMap[elemID]
		if !exists {
			fullData = &geodata.GetPointsDataByCategoryGroup{
				Measurements: make(map[int]*geodata.GroupMeasurement),
				IsCopied:     false, // Инициализация флага
			}
			fullDataMap[elemID] = fullData
		}

		measurement, exists := fullData.Measurements[mgID]
		if !exists {
			mg, exists := measurementGroupMap[mgID]
			if !exists {
				log.Printf("measurement group ID %d не найден в measurementGroupMap", mgID)
				continue
			}
			measurement = &geodata.GroupMeasurement{
				Name:      mg.Name,
				Unit:      mg.Unit,
				Rn:        mg.Rn,
				ZuluCoeff: getPointerIfValid(mg.ZuluCoeff),
				LersCoeff: getPointerIfValid(mg.LersCoeff),
				Data: geodata.GroupMeasurementsData{
					In:  mg.InLersVariable,
					Out: mg.OutLersVariable,
				},
			}
			fullData.Measurements[mgID] = measurement
		}

		lersVar := ""
		if lersVariable.Valid {
			lersVar = lersVariable.String
		}

		switch measurementRole {
		case "in":
			measurement.Data.In = lersVar
			if parsedVal != nil {
				measurement.CalculatedData.In = parsedVal
			} else {
				measurement.CalculatedData.In = nil
			}
		case "out":
			measurement.Data.Out = lersVar
			if parsedVal != nil {
				measurement.CalculatedData.Out = parsedVal
			} else {
				measurement.CalculatedData.Out = nil
			}
		default:
			log.Printf("неизвестная роль измерения '%s' для measurement_group_id %d и elem_id %d", measurementRole, mgID, elemID)
		}

		if insertedTS.Valid {
			if insertedTS.Time.Before(requestedTime) {
				fullData.IsCopied = true
			}
		} else {
			fullData.IsCopied = true
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("ошибка при итерации строк измерений: %v", err)
		return nil, errors.ErrInternalError
	}

	return fullDataMap, nil
}
