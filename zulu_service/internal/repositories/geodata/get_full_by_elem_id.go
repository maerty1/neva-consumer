package geodata

import (
	"context"
	"database/sql"
	"log"
	"time"

	"zulu_service/internal/config/errors"
	"zulu_service/internal/models/geodata"
)

func (r *repository) GetFullByElemID(ctx context.Context, elemID int, nDays int) (*geodata.FullElementData, error) {
	var fullData geodata.FullElementData

	var zwsTypeID int
	metaQuery := `
        SELECT
            em.title,
            em.address,
			em.zws_type_id
        FROM
            zulu.elems_metadata em
        WHERE
            em.elem_id = $1;
    `
	var title, address sql.NullString
	err := r.db.DB().QueryRow(ctx, metaQuery, elemID).Scan(&title, &address, &zwsTypeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFoundWithDetails("Element", elemID)
		}
		return nil, errors.NotFoundWithDetails("Объект", elemID)
	}
	if title.Valid {
		fullData.Title = title.String
	}

	if address.Valid {
		fullData.Address = address.String
	}

	measurementGroupsQuery := `
WITH cte_groups AS (
	WITH cte_category AS (
		SELECT
			oscf.full_category_id
		FROM
			public.measurement_categories mc
			JOIN object_state_configuration oscf ON oscf.full_category_id = mc.id
		WHERE
			oscf.zws_type_id = $1
	)
	SELECT
		group_id
	FROM
		public.category_to_group ctg
		JOIN cte_category ON cte_category.full_category_id = ctg.category_id
	WHERE
		ctg.category_id = cte_category.full_category_id
)
SELECT
	mg.id,
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
	join cte_groups on mg.id = cte_groups.group_id
	LEFT JOIN public.measurement_types in_mt ON mg.in = in_mt.id
	LEFT JOIN public.measurement_types out_mt ON mg.out = out_mt.id
	LEFT JOIN public.measurement_type_conversion mtc ON mtc.measurement_type_id = in_mt.id
	`

	measurementGroupsRows, err := r.db.DB().Query(ctx, measurementGroupsQuery, zwsTypeID)
	if err != nil {
		log.Printf("Ошибка выполнения запроса на получение групп измерений: %v", err)
		return nil, errors.ErrInternalError
	}
	defer measurementGroupsRows.Close()

	type measurementGroup struct {
		ID                   int
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

	measureQuery := `
        WITH max_ts AS (
            SELECT MAX(inserted_ts) AS latest_ts
            FROM zulu.object_records
            WHERE elem_id = $1
        )
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
			max_ts
			CROSS JOIN public.measurement_groups mg
			LEFT JOIN public.measurement_types mt ON mt.id = mg.in OR mt.id = mg.out
			LEFT JOIN zulu.object_records orc ON orc.parameter = mt.zulu_var
				AND orc.elem_id = $1
				AND orc.inserted_ts >= max_ts.latest_ts - ($2 * INTERVAL '1 day')
				AND orc.inserted_ts <= max_ts.latest_ts
		WHERE
			orc.td IS NULL
			AND mg.id = ANY ($3)
		ORDER BY
			orc.inserted_ts;
    `
	// Cp
	rows, err := r.db.DB().Query(ctx, measureQuery, elemID, nDays, relevantGroupIDs)
	if err != nil {
		log.Printf("Ошибка выполнения запроса на измерение: %v", err)
		return nil, errors.ErrInternalError
	}
	defer rows.Close()

	packetMap := make(map[string]map[int]*geodata.Measurement)

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
			return nil, errors.ErrInternalError
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
			packetMap[insertedTS] = make(map[int]*geodata.Measurement)
		}

		if _, exists := packetMap[insertedTS][mgID]; !exists {
			mg, exists := measurementGroupMap[mgID]
			if !exists {
				log.Printf("группа измерений ID %d не найдена в measurementGroupMap", mgID)
				continue
			}
			packetMap[insertedTS][mgID] = &geodata.Measurement{
				Name:      mg.Name,
				Unit:      mg.Unit,
				ZuluCoeff: getPointerIfValid(mg.ZuluCoeff),
				LersCoeff: getPointerIfValid(mg.LersCoeff),
				Data: geodata.MeasurementData{
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

	fullData.Packets = packetMap
	return &fullData, nil
}
