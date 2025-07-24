package reports

import (
	"context"
	"database/sql"
	"log"
	"zulu_service/internal/models/reports"
)

func (r *repository) GetBiDashboardFrameThree(ctx context.Context, elemID int) ([]reports.BiDashboardFrame, error) {
	var result []reports.BiDashboardFrame
	query := `
with ctp as (
            select elem_id 
            from zulu.zulu.objects_geometry_log i 
            where zws_type = 8)
SELECT
	q.elem_id block_id,
	COALESCE(em.address, em.title, q.elem_id::TEXT) block_name,
	q.qsum,
	q.specific_qsum*1000 specific_qsum
FROM
	public.qsum_by_branch q
	JOIN zulu.zulu.elems_metadata em ON em.elem_id = q.elem_id
	left join ctp on q.elem_id = ctp.elem_id
WHERE
	q.entrance_elem_id = $1::integer
	and ctp.elem_id is null
	AND qsum IS NOT NULL;
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
		var specificQsum sql.NullFloat64
		if err := rows.Scan(&block.BlockID, &blockName, &blockQsum, &specificQsum); err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			return nil, err
		}
		if blockName.Valid {
			block.BlockName = blockName.String
		}
		if blockQsum.Valid {
			block.Qsum = blockQsum.Float64
		}
		if specificQsum.Valid {
			block.SpecificQsum = specificQsum.Float64
		}

		result = append(result, block)
	}
	return result, nil
}
