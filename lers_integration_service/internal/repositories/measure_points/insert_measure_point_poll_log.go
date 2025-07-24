package measure_points

import (
	"context"
	"fmt"
	"log"
)

func (r *repository) InsertMeasurePointPollLog(ctx context.Context, pollID int, measurePointID int, accountID int, message string) error {
	log.Println("Вставка measure point log. pollID:", pollID, "message:", message)
	query := `
    INSERT INTO measure_points_poll_log (poll_id, message, measure_point_id, account_id, status)
	VALUES ($1, $2, $3, $4, $5)`

	if pollID != 0 {
		_, err := r.db.DB().Exec(ctx, query, pollID, message, measurePointID, accountID, "SUCCESS")
		if err != nil {
			return fmt.Errorf("не удалось вставить лог поллинга точки измерений: %v", err)
		}
		return nil
	}
	return nil

}
