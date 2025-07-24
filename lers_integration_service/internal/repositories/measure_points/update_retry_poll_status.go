package measure_points

import (
	"context"
	"fmt"
)

func (r *repository) UpdateRetryPollStatus(ctx context.Context, retryPoll RetryPollSessions, status string) error {
	tx, err := r.db.DB().Begin(ctx)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %v", err)
	}
	defer tx.Rollback(ctx)

	query := `
    UPDATE
		public.measure_points_poll_log
	SET
		status = $1
	WHERE
		poll_id = $2`

	query2 := `
	  UPDATE
		public.measure_points_poll_retry
	SET
		status = $1
	WHERE
		retrying_poll_id = $2
	`

	if _, err := tx.Exec(ctx, query, status, retryPoll.OriginalPollID); err != nil {
		return fmt.Errorf("не удалось обновить статус в measure_points_poll_log: %v", err)
	}

	if _, err := tx.Exec(ctx, query2, status, retryPoll.RetryPollID); err != nil {
		return fmt.Errorf("не удалось обновить статус в measure_points_poll_retry: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию: %v", err)
	}

	return nil
}
