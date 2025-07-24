package geodata

import (
	"context"
	"database/sql"

	"zulu_service/internal/config/errors"
)

func (r *repository) GetIconIdByElemId(ctx context.Context, elemID int) (int, error) { //string, error
	var iconID int

	metaQuery := `
	SELECT 
    	i.icon_id
	FROM 
    	public.icons i
	INNER JOIN 
 		public.point_to_icon p
	ON
    	p.icon_id = i.icon_id
	INNER JOIN
		zulu.elems_metadata m
	ON
		m.elem_id=p.elem_id
	WHERE 
    	zulu.elems_metadata.elem_id='$1'
	`

	err := r.db.DB().QueryRow(ctx, metaQuery, elemID).Scan(&iconID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.NotFoundWithDetails("Иконка", elemID)
		}
		return 0, errors.NotFoundWithDetails("Иконка", elemID)
	}
	return iconID, nil
}
