package geodata

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"zulu_service/internal/config/errors"
	"zulu_service/internal/models/geodata"
)

func (r *repository) GetPointsDataByCategoryGroup(ctx context.Context, elemIDs []int, categoryID int) (map[int]*geodata.GetPointsDataByCategoryGroup, error) {
	fullDataMap := make(map[int]*geodata.GetPointsDataByCategoryGroup)

	for _, elemID := range elemIDs {
		fullDataMap[elemID] = &geodata.GetPointsDataByCategoryGroup{
			Measurements: make(map[int]*geodata.GroupMeasurement),
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

	rows, err := r.db.DB().Query(ctx, measurementGroupsQuery, categoryID)
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
		ZuluCoeff            sql.NullFloat64
		LersCoeff            sql.NullFloat64
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
		fmt.Println(mg.ZuluCoeff.Valid)
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
    	with recursive scan as (
			(
			select
				concat(orc.elem_id, mt.id) val_id,
				mg.id AS measurement_group_id,
				orc.elem_id,
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
				AND orc.elem_id = ANY($1)
				AND orc.fd::date <= (now() at time zone 'utc')::date
				--AND orc.fd::date <= ($2::timestamp)::date
			WHERE
				mg.id = ANY($2)
			order by concat(orc.elem_id, mt.id) asc, fd::date desc, record_priority asc, fd desc
			limit 1 
		)
		union all (
			select r.*
			from scan 
			join lateral (
				select
				concat(orc.elem_id, mt.id) val_id,
				mg.id AS measurement_group_id,
				orc.elem_id,
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
				AND orc.elem_id = ANY($1)
				AND orc.fd::date <= (now() at time zone 'utc')::date
				--AND orc.fd::date <= ($2::timestamp)::date
			WHERE
				mg.id = ANY($2)
				and concat(orc.elem_id, mt.id) > scan.val_id
			order by concat(orc.elem_id, mt.id) asc, fd::date desc, record_priority asc, fd desc
			limit 1 
			) r 
			on true
		)
		)
		select 
			elem_id,
			measurement_group_id,
			measurement_group_name,
			measurement_unit,
			measurement_type_id,
			measurement_type_name,
			measurement_type_unit,
			measurement_role,
			val,
			rest_var,
			zulu_var 
		from scan
		;
    `

	rows, err = r.db.DB().Query(ctx, measureQuery, elemIDs, relevantGroupIDs)
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
				ZuluCoeff: getPointerIfValid(mg.ZuluCoeff),
				LersCoeff: getPointerIfValid(mg.LersCoeff),
				Rn:        mg.Rn,
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
	}

	if err = rows.Err(); err != nil {
		log.Printf("ошибка при итерации строк измерений.: %v", err)
		return nil, errors.ErrInternalError
	}

	return fullDataMap, nil
}
