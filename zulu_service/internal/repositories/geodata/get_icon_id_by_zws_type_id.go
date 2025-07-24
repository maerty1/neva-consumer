package geodata

import (
	"context"
	"database/sql"

	"zulu_service/internal/config/errors"

)

func (r *repository) GetIconIdByZwsType(ctx context.Context, elemID int) (int, error) { //string, error
	var iconID int

	metaQuery := `
	SELECT 
    	i.icon_id
	FROM 
    	public.icons_on_zws_type_id i
	INNER JOIN 
 		zulu.elems_metadata m
	ON
    	i.elem_id = m.elem_id
	WHERE 
    	zulu.elems_metadata.zws_type_id='$1'
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
