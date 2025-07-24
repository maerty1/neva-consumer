package geodata

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"zulu_service/internal/config/errors"
	"zulu_service/internal/models/geodata"
)

func (r *repository) GetPointDataByCategoryGroup(ctx context.Context, elemID int, categoryID int, timestamp string, nDays int) (*geodata.GetPointDataByCategoryGroup, error) {
	var fullData geodata.GetPointDataByCategoryGroup

	measurementGroupsQuery := `
    WITH cte_groups AS (
		SELECT
			group_id,
			rn
		FROM
			public.category_to_group ctg
		WHERE
			ctg.category_id = $1
		)
	SELECT
		mg.id,
		cte_groups.rn,
		mg.group_front_desc,
		mg.measurement_unit,
		mg.in,
		mg.out,
		in_mt.rest_var AS in_lers_var,
		out_mt.rest_var AS out_lers_var,
		mtc.lers_coeff,
		mtc.zulu_coeff
	FROM
		public.measurement_groups mg
		JOIN cte_groups ON mg.id = cte_groups.group_id
		LEFT JOIN public.measurement_types in_mt ON mg.in = in_mt.id
		LEFT JOIN public.measurement_types out_mt ON mg.out = out_mt.id
		LEFT JOIN public.measurement_type_conversion mtc ON mtc.measurement_type_id = in_mt.id
	`

	measurementGroupsRows, err := r.db.DB().Query(ctx, measurementGroupsQuery, categoryID)
	if err != nil {
		log.Printf("Ошибка выполнения запроса на получение групп измерений: %v", err)
		return nil, errors.ErrInternalError
	}
	defer measurementGroupsRows.Close()

	type measurementGroup struct {
		ID                   int
		Rn                   int
		Name                 string
		Unit                 string
		InMeasurementTypeID  int
		OutMeasurementTypeID int
		InLersVariable       string
		OutLersVariable      string
		ZuluCoeff            sql.NullFloat64
		LersCoeff            sql.NullFloat64
	}

	var measurementGroups []measurementGroup
	var relevantGroupIDs []int
	measurementGroupMap := make(map[int]measurementGroup)

	for measurementGroupsRows.Next() {
		var mg measurementGroup
		err := measurementGroupsRows.Scan(
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
			log.Printf("Ошибка сканирования строки групп измерений: %v", err)
			return nil, errors.ErrInternalError
		}
		measurementGroups = append(measurementGroups, mg)
		measurementGroupMap[mg.ID] = mg
		relevantGroupIDs = append(relevantGroupIDs, mg.ID)
	}

	if err = measurementGroupsRows.Err(); err != nil {
		log.Printf("Ошибка итерации строк групп измерений: %v", err)
		return nil, errors.ErrInternalError
	}

	// Определяем SQL-запрос и параметры на основе переданных параметров
	var measureQuery string
	var queryParams []interface{}

	if timestamp != "" {
		// Запрос для получения данных по конкретному timestamp
		measureQuery = `
	with recursive scan as (
		(
		SELECT
			orc.inserted_ts,
			mg.id AS measurement_group_id,
			mg.group_front_desc AS measurement_group_name,
			mg.measurement_unit AS measurement_unit,
			mt.id AS measurement_type_id,
			mt.front_desc AS measurement_type_name,
			mt.zulu_un AS measurement_type_unit,
			CASE
				WHEN mt.id = mg.in THEN 'in'
				WHEN mt.id = mg.out THEN 'out'
				ELSE 'unknown'
			END AS measurement_role,
			orc.val,
			mt.rest_var,
			mt.zulu_var
		FROM
			public.measurement_groups mg
		JOIN 
			public.measurement_types mt ON mt.id = mg.in OR mt.id = mg.out
		LEFT JOIN 
			zulu.object_records orc ON orc.parameter = mt.zulu_var
			AND orc.elem_id = $1
			AND orc.fd::date <= ($2::timestamp)::date
		WHERE
			mg.id = ANY ($3)
		order by mt.id asc, fd::date desc, record_priority asc, fd desc
		limit 1 
	)
		union all (
		select r.*
		from scan 
		join lateral (
			SELECT
				orc.inserted_ts,
				mg.id AS measurement_group_id,
				mg.group_front_desc AS measurement_group_name,
				mg.measurement_unit AS measurement_unit,
				mt.id AS measurement_type_id,
				mt.front_desc AS measurement_type_name,
				mt.zulu_un AS measurement_type_unit,
				CASE
					WHEN mt.id = mg.in THEN 'in'
					WHEN mt.id = mg.out THEN 'out'
					ELSE 'unknown'
				END AS measurement_role,
				orc.val,
				mt.rest_var,
				mt.zulu_var
			FROM
				public.measurement_groups mg
			JOIN 
				public.measurement_types mt ON mt.id = mg.in OR mt.id = mg.out
			LEFT JOIN 
				zulu.object_records orc ON orc.parameter = mt.zulu_var
				AND orc.elem_id = $1
				AND orc.fd::date <= ($2::timestamp)::date
			WHERE
				mg.id = ANY ($3)
				and mt.id > scan.measurement_type_id
			order by mt.id asc, fd::date desc, record_priority asc, fd desc
			limit 1 
		) r 
		on true
	)
	)
	select * from scan
	;
	`
		fmt.Println(relevantGroupIDs)
		queryParams = []interface{}{elemID, timestamp, relevantGroupIDs}
	} else {
		// Запрос для получения данных за последние n дней
		measureQuery = `
	with recursive scan as (
		(
		SELECT
			orc.inserted_ts,
			mg.id AS measurement_group_id,
			mg.group_front_desc AS measurement_group_name,
			mg.measurement_unit AS measurement_unit,
			mt.id AS measurement_type_id,
			mt.front_desc AS measurement_type_name,
			mt.zulu_un AS measurement_type_unit,
			CASE
				WHEN mt.id = mg.in THEN 'in'
				WHEN mt.id = mg.out THEN 'out'
				ELSE 'unknown'
			END AS measurement_role,
			orc.val,
			mt.rest_var,
			mt.zulu_var
		FROM
			public.measurement_groups mg
		JOIN 
			public.measurement_types mt ON mt.id = mg.in OR mt.id = mg.out
		LEFT JOIN 
			zulu.object_records orc ON orc.parameter = mt.zulu_var
			AND orc.elem_id = $1
			AND orc.fd::date <= ($2::timestamp)::date
		WHERE
			mg.id = ANY ($3)
		order by mt.id asc, record_priority asc, fd desc
		limit 1 
	)
		union all (
		select r.*
		from scan 
		join lateral (
			SELECT
				orc.inserted_ts,
				mg.id AS measurement_group_id,
				mg.group_front_desc AS measurement_group_name,
				mg.measurement_unit AS measurement_unit,
				mt.id AS measurement_type_id,
				mt.front_desc AS measurement_type_name,
				mt.zulu_un AS measurement_type_unit,
				CASE
					WHEN mt.id = mg.in THEN 'in'
					WHEN mt.id = mg.out THEN 'out'
					ELSE 'unknown'
				END AS measurement_role,
				orc.val,
				mt.rest_var,
				mt.zulu_var
			FROM
				public.measurement_groups mg
			JOIN 
				public.measurement_types mt ON mt.id = mg.in OR mt.id = mg.out
			LEFT JOIN 
				zulu.object_records orc ON orc.parameter = mt.zulu_var
				AND orc.elem_id = $1
				AND orc.fd::date <= (now() at time zone 'utc')::date
			WHERE
				mg.id = ANY ($3)
				and mt.id > scan.measurement_type_id
			order by mt.id asc, record_priority asc, fd desc
			limit 1 
		) r 
		on true
	)
	)
	select * from scan
	;
	`
		queryParams = []interface{}{elemID, nDays, relevantGroupIDs}
	}

	rows, err := r.db.DB().Query(ctx, measureQuery, queryParams...)
	if err != nil {
		log.Printf("Ошибка выполнения запроса на измерение: %v", err)
		return nil, errors.ErrInternalError
	}
	defer rows.Close()

	packetMap := make(map[string]map[int]*geodata.GroupMeasurement)

	for rows.Next() {
		var (
			timeInsertedTS  sql.NullTime
			mgID            int
			mgName          string
			mgUnit          string
			mtID            int
			mtName          sql.NullString
			mtUnit          sql.NullString
			measurementRole string
			val             sql.NullFloat64 // orc.val
			lersVariable    sql.NullString  // mt.rest_var
			zuluVariable    sql.NullString  // mt.zulu_var
		)

		err := rows.Scan(
			&timeInsertedTS,
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
		)
		if err != nil {
			log.Printf("Ошибка сканирования строки измерения: %v", err)
			continue
			// return nil, errors.ErrInternalError
		}

		var insertedTS string
		if timeInsertedTS.Valid {
			insertedTS = timeInsertedTS.Time.Format("2006-01-02")
		} else {
			insertedTS = time.Now().Format("2006-01-02")
		}

		var lersVar string
		if lersVariable.Valid {
			lersVar = lersVariable.String
		}

		if _, exists := packetMap[insertedTS]; !exists {
			packetMap[insertedTS] = make(map[int]*geodata.GroupMeasurement)
		}

		if _, exists := packetMap[insertedTS][mgID]; !exists {
			mg, exists := measurementGroupMap[mgID]
			if !exists {
				log.Printf("группа измерений ID %d не найдена в measurementGroupMap", mgID)
				continue
			}
			packetMap[insertedTS][mgID] = &geodata.GroupMeasurement{
				Name:      mg.Name,
				Unit:      mg.Unit,
				ZuluCoeff: getPointerIfValid(mg.ZuluCoeff),
				LersCoeff: getPointerIfValid(mg.LersCoeff),
				Rn:        mg.Rn,
				Data: geodata.GroupMeasurementsData{
					In:  mg.InLersVariable,
					Out: mg.OutLersVariable,
				},
			}
		}

		md := packetMap[insertedTS][mgID]

		if measurementRole == "in" {
			md.Data.In = lersVar
			if val.Valid {
				md.CalculatedData.In = &val.Float64
			} else {
				md.CalculatedData.In = nil
			}
		} else if measurementRole == "out" {
			md.Data.Out = lersVar
			if val.Valid {
				md.CalculatedData.Out = &val.Float64
			} else {
				md.CalculatedData.Out = nil
			}
		} else {
			log.Printf("Неизвестная роль для measure_group_id %d at %v", mgID, insertedTS)
			continue
		}
	}

	fullData.Measurements = packetMap
	return &fullData, nil
}
