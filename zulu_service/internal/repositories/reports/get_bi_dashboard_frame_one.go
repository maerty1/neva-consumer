package reports

import (
	"context"
	"log"
	"zulu_service/internal/models/reports"
)

func (r *repository) GetBiDashboardFrame(ctx context.Context) ([]reports.BiDashboardFrame, error) {
	var result []reports.BiDashboardFrame
	query := `
		with recursive scan as (
		    (
		    SELECT
			r.elem_id block_id,
			COALESCE(em.title, r.elem_id::TEXT) block_name,
			val::NUMERIC*24 qsum
		FROM
			zulu.zulu.object_records r
			LEFT JOIN zulu.zulu.elems_metadata em ON em.elem_id = r.elem_id
		WHERE
			r.elem_id IN(
				SELECT
					elem_id
				FROM
					zulu.zulu.objects_geometry_log i
				WHERE
					zws_type = 1
			)
			AND "parameter" IN ('Qsum')
		    order by r.elem_id asc, fd desc
		    limit 1 
		  )
		  union all (
		    select r.*
		    from scan 
		    join lateral (
		            SELECT
			r.elem_id block_id,
			COALESCE(em.title, r.elem_id::TEXT) block_name,
			val::NUMERIC*24 qsum
		FROM
			zulu.zulu.object_records r
			LEFT JOIN zulu.zulu.elems_metadata em ON em.elem_id = r.elem_id
		WHERE
			r.elem_id IN(
				SELECT
					elem_id
				FROM
					zulu.zulu.objects_geometry_log i
				WHERE
					zws_type = 1
			)
			AND "parameter" IN ('Qsum')
			and r.elem_id > scan.block_id
		    order by r.elem_id asc, fd desc
		    limit 1 
		    ) r 
		    on true
		 )
		)
		select * from scan
		;
    `

	rows, err := r.db.DB().Query(ctx, query)
	if err != nil {
		log.Printf("Ошибка при получении состояний объектов: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var block reports.BiDashboardFrame
		if err := rows.Scan(&block.BlockID, &block.BlockName, &block.Qsum); err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			return nil, err
		}
		result = append(result, block)
	}
	return result, nil
}
