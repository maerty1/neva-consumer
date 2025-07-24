package zulu

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"zulu_updater/internal/models"
)

func (r repository) InsertRecords(ctx context.Context, fields []models.Field, elemId int) error {
	insertedTs := time.Now().Format("2006-01-02 15:04:05.00")

	var values []string
	for _, field := range fields {
		values = append(values, fmt.Sprintf("('%d', '%s', '%s', '%s')", elemId, field.Name, insertedTs, field.Value))
	}

	query := fmt.Sprintf(`
		INSERT INTO zulu.object_records_new (elem_id, val_name, inserted_ts, val)
		VALUES %s
	`, strings.Join(values, ","))

	row, err := r.db.DB().Query(ctx, query)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	var temp int
	var res []int
	for row.Next() {
		err = row.Scan(&temp)
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		res = append(res, temp)
	}
	return nil
}
