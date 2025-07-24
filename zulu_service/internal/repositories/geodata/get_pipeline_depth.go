package geodata

import (
	"context"
	"log"
	"zulu_service/internal/config/errors"
	"zulu_service/internal/models/geodata"
)

func (r *repository) GetPipelineDepth(ctx context.Context) (map[int]geodata.PipelineDepth, error) {
	pipelineDepth := make(map[int]geodata.PipelineDepth)

	query := `
WITH border_objects AS (
	SELECT
		i.elem_id,
		ST_StartPoint(i.zws_geometry) start_point,
		MAX(i2.elem_id) start_elem_id,
		ST_EndPoint(i.zws_geometry) end_point,
		MAX(i3.elem_id) end_elem_id
	FROM
		zulu.zulu.objects_geometry_log i
		LEFT JOIN zulu.zulu.objects_geometry_log i2 ON i2.zws_type != 6
		AND ST_StartPoint(i.zws_geometry) = i2.zws_geometry
		LEFT JOIN zulu.zulu.objects_geometry_log i3 ON i3.zws_type != 6
		AND ST_EndPoint(i.zws_geometry) = i3.zws_geometry
	WHERE
		i.zws_type = 6
	GROUP BY
		i.elem_id,
		i.zws_geometry
)
SELECT
	b.elem_id,
	-- start_point,
	r.val::NUMERIC start_h_geo,
	-- end_point,
	r2.val::NUMERIC end_h_geo,
	(r.val::NUMERIC + r2.val::NUMERIC)/ 2 avg_h_geo
FROM
	border_objects b
	JOIN zulu.zulu.object_records r ON start_elem_id = r.elem_id
	AND r."parameter" = 'H_geo'
	JOIN zulu.zulu.object_records r2 ON end_elem_id = r2.elem_id
	AND r2."parameter" = 'H_geo';
	`

	rows, err := r.db.DB().Query(ctx, query)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return nil, errors.ErrInternalError
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var startHgeo, endHGeo, avgHgeo float64

		if err := rows.Scan(&id, &startHgeo, &endHGeo, &avgHgeo); err != nil {
			log.Printf("Ошибка сканирования строки: %v", err)
			return nil, errors.ErrInternalError
		}

		pipelineDepth[id] = geodata.PipelineDepth{
			StartHgeo: startHgeo,
			EndHgeo:   endHGeo,
			AvgHgeo:   avgHgeo,
		}

	}
	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при итерации по строкам: %v", err)
		return nil, errors.ErrInternalError
	}

	return pipelineDepth, nil

}
