package measure_points

import (
	"context"
	"lers_integration_service/internal/db"
	"lers_integration_service/internal/models"
)

type Repository interface {
	FindAccountsToSync(ctx context.Context) ([]models.AccountToSync, error)
	FindPollSessionsToRetry(ctx context.Context, accountID int) ([]PollSessionsToRetry, error)
	FindRetryPollSessions(ctx context.Context, accountID int) ([]RetryPollSessions, error)
	// TODO: Придумать название получше
	FindPollSessionsToRetry2(ctx context.Context, accountID int) ([]PollSessionsToRetry2, error)

	GetLastMeasurePointDatetime(ctx context.Context, measurePointID int) (string, error)

	InsertMeasurePointData(ctx context.Context, measurePointID int, datetime string, values string) error
	InsertMeasurePointDayData(ctx context.Context, measurePointID int, datetime string, values string) error
	InsertMeasurePoint(ctx context.Context, accountID int, measurePointID int, deviceID int, title string, fullTitle string, address string, system_type string) error
	InsertMeasurePointPollRetry(ctx context.Context, originalPollID int, retryPollID int, status string) error
	InsertSyncLog(ctx context.Context, accountID int, measurePointID int, level, message string) error
	InsertMeasurePointPollLog(ctx context.Context, pollID int, measurePointID int, accountID int, message string) error
	InsertMeasurePointDayDataBatch(ctx context.Context, data []MeasurePointsDataDay) error

	UpdatePollStatus(ctx context.Context, pollID int, status string) error
	UpdateRetryPollStatus(ctx context.Context, retryPoll RetryPollSessions, status string) error
}

var _ Repository = (*repository)(nil)

type repository struct {
	db db.PostgresClient
}

func NewRepository(db db.PostgresClient) *repository {
	return &repository{
		db: db,
	}
}
