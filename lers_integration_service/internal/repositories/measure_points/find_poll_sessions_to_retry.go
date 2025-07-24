package measure_points

import (
	"context"
)

type PollSessionsToRetry struct {
	PollID         int
	MeasurePointID int
}

func (r *repository) FindPollSessionsToRetry(ctx context.Context, accountID int) ([]PollSessionsToRetry, error) {
	query := `
	SELECT
		poll_id,
		measure_point_id
	FROM
		measure_points_poll_log
	WHERE
		status != 'SUCCESS'
		AND account_id = $1;
	`

	rows, err := r.db.DB().Query(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pollSessions []PollSessionsToRetry
	for rows.Next() {
		var session PollSessionsToRetry
		err = rows.Scan(&session.PollID, &session.MeasurePointID)
		if err != nil {
			return nil, err
		}

		pollSessions = append(pollSessions, session)
	}

	return pollSessions, nil
}
