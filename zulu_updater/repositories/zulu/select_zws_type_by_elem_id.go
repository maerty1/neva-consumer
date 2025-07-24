package zulu

import (
	"context"
)

func (r repository) SelectZwsTypeByElemId(ctx context.Context, elemId int) (int, error) {
	query := `SELECT zws_type FROM zulu.objects_geometry_log WHERE elem_id = $1 LIMIT 1;`
	row := r.db.DB().QueryRow(ctx, query, elemId)
	var res int
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}
