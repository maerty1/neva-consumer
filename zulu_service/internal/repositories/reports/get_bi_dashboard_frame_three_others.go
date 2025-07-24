package reports

import (
	"context"
	"database/sql"
	"log"
	"zulu_service/internal/models/reports"
)

func (r *repository) GetBiDashboardFrameThreeOthers(ctx context.Context, elemID int) ([]reports.BiDashboardFrame, error) {
	var result []reports.BiDashboardFrame
	query := `
with ctp as (
            select elem_id 
            from zulu.zulu.objects_geometry_log i 
            where zws_type = 8
)
select q.elem_id 			            block_id, 
	   coalesce(em.title, q.elem_id::text)    block_name,
	   qsum
from public.qsum_by_branch q
left join ctp 
on entrance_elem_id = ctp.elem_id
or q.elem_id = ctp.elem_id
left join zulu.zulu.elems_metadata em 
on em.elem_id = q.elem_id 
where ctp.elem_id is null
and istok_id = $1::integer;
    `

	rows, err := r.db.DB().Query(ctx, query, elemID)
	if err != nil {
		log.Printf("Ошибка при получении состояний объектов: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var block reports.BiDashboardFrame
		var blockName sql.NullString
		var blockQsum sql.NullFloat64

		if err := rows.Scan(&block.BlockID, &blockName, &blockQsum); err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			return nil, err
		}
		if blockName.Valid {
			block.BlockName = blockName.String
		}
		if blockQsum.Valid {
			block.Qsum = blockQsum.Float64
		}
		result = append(result, block)
	}
	return result, nil
}
