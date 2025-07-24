package geodata

import (
	"context"
	"database/sql"

	"zulu_service/internal/config/errors"
)

func (r *repository) GetSchemaIdByElemId(ctx context.Context, elemID int) (int, error) {
	var schemaID int

	metaQuery := `
    SELECT
        em.schema_id
    FROM
        zulu.elems_metadata em
    WHERE
        em.elem_id = $1;
    `
	err := r.db.DB().QueryRow(ctx, metaQuery, elemID).Scan(&schemaID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.NotFoundWithDetails("Схема", elemID)
		}
		return 0, errors.NotFoundWithDetails("Схема", elemID)
	}
	return schemaID, nil
}
