package measure_points

import (
	"context"
	"log"
)

func (r *repository) InsertMeasurePoint(ctx context.Context, accountID int, measurePointID int, deviceID int, title string, fullTitle string, address string, system_type string) error {
	query := `
	INSERT INTO measure_points (id, account_id, title, device_id, full_title, address, system_type)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (id) DO UPDATE SET title = $3, device_id = $4, full_title = $5, address = $6, system_type=$7`
	_, err := r.db.DB().Exec(ctx, query, measurePointID, accountID, title, deviceID, fullTitle, address, system_type)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println("Успешная вставка")
	return nil
}
