package measure_points

import (
	"context"
)

type RetryPollSessions struct {
	OriginalPollID int
	RetryPollID    int
}

func (r *repository) FindRetryPollSessions(ctx context.Context, accountID int) ([]RetryPollSessions, error) {
	query := `
	SELECT 
		r.original_poll_id, 
		r.retrying_poll_id
	FROM 
		public.measure_points_poll_retry r
	JOIN 
		public.measure_points_poll_log l 
	ON 
		r.original_poll_id = l.poll_id
	WHERE 
		(r.status != 'SUCCESS' AND r.status != 'FAILED') 
		AND l.account_id = $1;
	`

	rows, err := r.db.DB().Query(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pollSessions []RetryPollSessions
	for rows.Next() {
		var session RetryPollSessions
		err = rows.Scan(&session.OriginalPollID, &session.RetryPollID)
		if err != nil {
			return nil, err
		}

		pollSessions = append(pollSessions, session)
	}

	return pollSessions, nil
}
