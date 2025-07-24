package reports

import (
	"context"
	"database/sql"
	"log"
	"zulu_service/internal/models/reports"
)

func (r *repository) GetBiDashboardFrameTwo(ctx context.Context, elemID int) ([]reports.BiDashboardFrame, error) {
	var result []reports.BiDashboardFrame
	query := `
with branches as (
            select distinct entrance_elem_id branch_id from public.qsum_by_branch
where istok_id = $1
),
src as (
select 
            branch_id
 			block_id, 
	        entrance_pipe_id, 
	        q.elem_id,
	        qsum
        from public.qsum_by_branch q
        join branches 
        on entrance_elem_id = branch_id
        and istok_id = $1     
        where q.elem_id not in (select * from branches)
        )
select 
    block_id, 
    case 
        when block_id = $1
        then coalesce (b.branch_name, 'Остальные')
        else coalesce (b.branch_name, em.title, block_id::text)
    end         block_name, 
    SUM(qsum)   qsum
from src
left join zulu.zulu.elems_metadata em 
on em.elem_id = block_id
left join zulu.zulu.branch_names b
on b.entrance_elem_id = block_id
group by block_id, em.title, branch_name;
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

		if err := rows.Scan(&block.BlockID, &blockName, &block.Qsum); err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			return nil, err
		}
		if blockName.Valid {
			block.BlockName = blockName.String
		}
		result = append(result, block)
	}
	return result, nil
}
