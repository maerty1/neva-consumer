package measure_points

import (
	"context"
	"fmt"
	"log"
)

func (r *repository) InsertMeasurePointPollRetry(ctx context.Context, originalPollID int, retryPollID int, status string) error {
	log.Println("Insert measure point retry", retryPollID, "for", originalPollID)
	query := `
    INSERT INTO measure_points_poll_retry (original_poll_id, retrying_poll_id, status)
	VALUES ($1, $2, $3)`
	_, err := r.db.DB().Exec(ctx, query, originalPollID, retryPollID, status)
	if err != nil {
		return fmt.Errorf("не удалось вставить лог поллинга точки измерений: %v", err)
	}
	return nil
}
