package geodata

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"zulu_service/internal/config/errors"
	"zulu_service/internal/models/geodata"
)

func (r *repository) GetFilteredPoints(ctx context.Context, elementIDs []int, zwsTypeIDs []int, timestamp string) ([]geodata.Point, error) {
	var pointsData []geodata.Point

	var requestedTime time.Time
	var err error

	requestedTime, err = time.Parse(time.RFC3339, timestamp)
	if err != nil {
		layout := "2006-01-02"
		requestedTime, err = time.Parse(layout, timestamp)
		if err != nil {
			log.Printf("Неверный формат timestamp: %v", err)
			// return nil, errors.ErrInvalidInput

			requestedTime = time.Now()
		}

	}

	query := `
SELECT
    em.elem_id,
    em.title,
    em.address,
    CASE
        WHEN ST_GeometryType(ogl.zws_geometry) = 'ST_Point' THEN 
            json_build_array(ST_Y(ogl.zws_geometry), ST_X(ogl.zws_geometry))
        ELSE 
            NULL
    END AS coordinates,
    osc.collapsed_category_id,
    cg.group_id,
    mg.id AS measurement_group_id,
    it_in.rest_var AS in_variable,
    it_out.rest_var AS out_variable,
    em.zws_type_id AS type,
	mtc.lers_coeff,
	orr.inserted_ts
FROM
	zulu.elems_metadata em
	LEFT JOIN zulu.objects_geometry_log ogl ON em.elem_id = ogl.elem_id
	LEFT JOIN public.object_state_configuration osc ON em.zws_type_id = osc.zws_type_id
	LEFT JOIN public.category_to_group cg ON osc.collapsed_category_id = cg.category_id
	LEFT JOIN public.measurement_groups mg ON cg.group_id = mg.id
	LEFT JOIN public.measurement_types it_in ON mg.in = it_in.id
	LEFT JOIN public.measurement_types it_out ON mg.out = it_out.id
	LEFT JOIN public.measurement_type_conversion mtc ON mtc.measurement_type_id = it_in.id
	INNER JOIN zulu.object_records orr ON em.elem_id = orr.elem_id
	AND orr.is_deleted = FALSE
WHERE
    em.elem_id = ANY ($1) AND
	em.zws_type_id = ANY ($2)
ORDER BY
    em.elem_id;
    `

	rows, err := r.db.DB().Query(ctx, query, elementIDs, zwsTypeIDs)
	if err != nil {
		log.Printf("ошибка выполнения запроса: %v", err)
		return nil, errors.ErrInternalError
	}
	defer rows.Close()

	pointsMap := make(map[int]*geodata.Point)

	for rows.Next() {
		var (
			elemID              int
			title               sql.NullString
			address             sql.NullString
			coordinatesStr      sql.NullString
			collapsedCategoryID sql.NullInt32
			groupID             sql.NullInt32
			measurementGroupID  sql.NullInt32
			inVariable          sql.NullString
			outVariable         sql.NullString
			typeID              int
			lersCoeff           sql.NullFloat64
			insertedTS          sql.NullTime
		)

		err := rows.Scan(
			&elemID,
			&title,
			&address,
			&coordinatesStr,
			&collapsedCategoryID,
			&groupID,
			&measurementGroupID,
			&inVariable,
			&outVariable,
			&typeID,
			&lersCoeff,
			&insertedTS,
		)
		if err != nil {
			log.Printf("ошибка сканирования строки: %v", err)
			return nil, errors.ErrInternalError
		}

		point, exists := pointsMap[elemID]
		if !exists {
			var coordinates []float64
			if coordinatesStr.Valid {
				if err := json.Unmarshal([]byte(coordinatesStr.String), &coordinates); err != nil {
					log.Printf("ошибка парсинга координат для elem_id %d: %v", elemID, err)
					return nil, errors.ErrInternalError
				}
			}

			point = &geodata.Point{
				ElemID:            elemID,
				Title:             &title.String,
				Address:           &address.String,
				MeasurementGroups: make(map[int]geodata.MeasurementGroup),
				Coordinates:       coordinates,
				HasAccident:       false,
				Type:              typeID,
			}
			pointsMap[elemID] = point
		}

		if groupID.Valid && measurementGroupID.Valid {
			mgID := int(measurementGroupID.Int32)
			if _, exists := point.MeasurementGroups[mgID]; !exists {
				point.MeasurementGroups[mgID] = geodata.MeasurementGroup{}
			}

			mg := point.MeasurementGroups[mgID]

			if inVariable.Valid {
				mg.In = inVariable.String
			}
			if outVariable.Valid {
				mg.Out = outVariable.String
			}
			if lersCoeff.Valid {
				mg.Coeff = &lersCoeff.Float64
			}

			point.MeasurementGroups[mgID] = mg
		}

		if insertedTS.Valid {
			if insertedTS.Time.Before(requestedTime) {
				point.IsCopied = true
			}
		} else {
			point.IsCopied = true
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("ошибка итерации по строкам: %v", err)
		return nil, errors.ErrInternalError
	}

	for _, point := range pointsMap {
		pointsData = append(pointsData, *point)
	}

	return pointsData, nil
}
