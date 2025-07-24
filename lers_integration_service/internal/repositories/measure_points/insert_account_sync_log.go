package measure_points

import (
	"context"
	"fmt"
)

func (r *repository) InsertSyncLog(ctx context.Context, accountID int, measurePointID int, level, message string) error {
	query := `
    INSERT INTO accounts_sync_log (account_id, measure_point_id, level, message)
	VALUES ($1, $2, $3, $4)`
	_, err := r.db.DB().Exec(ctx, query, accountID, measurePointID, level, message)
	if err != nil {
		return fmt.Errorf("не удалось вставить лог синхронизации: %v", err)
	}
	return nil
}
