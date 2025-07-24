package geodata

import (
	"context"
	"database/sql"
	"log"
	"zulu_service/internal/config/errors"
	"zulu_service/internal/models/geodata"
)

func (r *repository) GetPointCategories(ctx context.Context, elemID int) (geodata.PointWithCategories, error) {
	pointWithCategories := geodata.PointWithCategories{
		Categories: []geodata.Category{},
	}

	query := `
	SELECT
		em.address,
		em.title,
		mc.id,
		mc.name,
		mc.expanded_rows_qty,
		mc.is_open,
		mc.cut_type,
		em.zws_type_id
	FROM
		zulu.elems_metadata AS em
		JOIN public.object_state_to_measurement_category AS obj_state ON obj_state.zws_type_id = em.zws_type_id
		JOIN public.measurement_categories AS mc ON obj_state.measurement_category_id = mc.id
	WHERE
		elem_id = $1
	ORDER BY obj_state.rn
	`

	rows, err := r.db.DB().Query(ctx, query, elemID)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return pointWithCategories, errors.ErrInternalError
	}

	defer rows.Close()

	rowsProcessed := 0
	for rows.Next() {
		rowsProcessed++
		var categoryName, categoryType string
		var title, address sql.NullString
		var categoryID, zwsTypeID, categoryMaxValues int
		var categoryIsOpen int

		if err := rows.Scan(&address, &title, &categoryID, &categoryName, &categoryMaxValues, &categoryIsOpen, &categoryType, &zwsTypeID); err != nil {
			log.Printf("Ошибка сканирования строки: %v", err)
			return pointWithCategories, errors.ErrInternalError
		}

		if address.Valid {
			pointWithCategories.Address = &address.String

		}
		if title.Valid {
			pointWithCategories.Title = &title.String

		}
		pointWithCategories.Type = zwsTypeID
		pointWithCategories.Categories = append(pointWithCategories.Categories, geodata.Category{Name: categoryName, ID: categoryID, Type: categoryType, MaxValues: &categoryMaxValues, IsOpen: categoryIsOpen == 1})
	}

	if rowsProcessed == 0 {
		return geodata.PointWithCategories{}, errors.NotFoundWithDetails("объект", elemID)

	}

	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при итерации по строкам: %v", err)
		return pointWithCategories, errors.ErrInternalError
	}

	return pointWithCategories, nil
}
