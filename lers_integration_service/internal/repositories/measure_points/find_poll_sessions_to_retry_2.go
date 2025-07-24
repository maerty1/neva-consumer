package measure_points

import (
	"context"
)

type PollSessionsToRetry2 struct {
	PollID         int
	MeasurePointID int
}

func (r *repository) FindPollSessionsToRetry2(ctx context.Context, accountID int) ([]PollSessionsToRetry2, error) {
	// Обновляем статус старых записей на "OUTDATED" для каждого measure_point_id, только если их больше одной,
	// и последняя запись не имеет статус "SUCCESS".
	updateQuery := `
	WITH PollCounts AS (
		SELECT 
			measure_point_id, 
			COUNT(*) AS poll_count,
			MAX(poll_id) AS latest_poll_id
		FROM 
			measure_points_poll_log
		WHERE 
			account_id = $1
			AND status != 'OUTDATED'
		GROUP BY 
			measure_point_id
		HAVING 
			COUNT(*) > 1
	)
	UPDATE 
		measure_points_poll_log l
	SET 
		status = 'OUTDATED'
	FROM 
		PollCounts pc
	WHERE 
		l.measure_point_id = pc.measure_point_id
		AND l.poll_id != pc.latest_poll_id
		AND l.account_id = $1
		AND l.status != 'SUCCESS'
		AND pc.latest_poll_id NOT IN (
			SELECT poll_id 
			FROM measure_points_poll_log 
			WHERE status = 'SUCCESS'
		);
	`

	_, err := r.db.DB().Exec(ctx, updateQuery, accountID)
	if err != nil {
		return nil, err
	}

	// Основной запрос, который возвращает актуальные poll сессии для ретрая
	query := `
	WITH RetryCounts AS (
		SELECT 
			original_poll_id, 
			COUNT(*) AS retry_count
		FROM 
			measure_points_poll_retry
		WHERE 
			DATE(created_at) = CURRENT_DATE
		GROUP BY 
			original_poll_id
	)
	SELECT
		l.poll_id,
		l.measure_point_id
	FROM
		measure_points_poll_log l
	LEFT JOIN 
		RetryCounts r
	ON 
		l.poll_id = r.original_poll_id
	WHERE
		l.status != 'SUCCESS'
		AND l.status != 'OUTDATED'
		AND l.account_id = $1
		AND (r.retry_count IS NULL OR r.retry_count < 3);
	`

	rows, err := r.db.DB().Query(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pollSessions []PollSessionsToRetry2
	for rows.Next() {
		var session PollSessionsToRetry2
		err = rows.Scan(&session.PollID, &session.MeasurePointID)
		if err != nil {
			return nil, err
		}

		pollSessions = append(pollSessions, session)
	}

	return pollSessions, nil
}
