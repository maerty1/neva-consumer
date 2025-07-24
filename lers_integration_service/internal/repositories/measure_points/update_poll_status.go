package measure_points

import (
	"context"
	"fmt"
)

func (r *repository) UpdatePollStatus(ctx context.Context, pollID int, status string) error {
	query := `
    UPDATE
		public.measure_points_poll_log
	SET
		status = $1
	WHERE
		poll_id = $2`
	_, err := r.db.DB().Exec(ctx, query, status, pollID)
	if err != nil {
		return fmt.Errorf("не удалось вставить лог синхронизации: %v", err)
	}
	return nil
}
