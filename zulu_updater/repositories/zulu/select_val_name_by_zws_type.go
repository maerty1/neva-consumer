package zulu

import (
	"context"
)

func (r repository) SelectValNameByZwsType(ctx context.Context, zwsType int, extractionType string) ([]string, error) {
	query := `select fields from zulu.extraction_parameters where extraction_type = $1 and zws_type = $2`
	row := r.db.DB().QueryRow(ctx, query, extractionType, zwsType)
	var res []string
	err := row.Scan(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
