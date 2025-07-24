package reports

import (
	"context"

	"zulu_service/internal/db"
	"zulu_service/internal/models/reports"
)

type Repository interface {
	GetBiDashboardFrame(ctx context.Context) ([]reports.BiDashboardFrame, error)
	GetBiDashboardFrameTwo(ctx context.Context, elemID int) ([]reports.BiDashboardFrame, error)
	GetBiDashboardFrameThree(ctx context.Context, elemID int) ([]reports.BiDashboardFrame, error)
	GetBiDashboardFrameThreeOthers(ctx context.Context, elemID int) ([]reports.BiDashboardFrame, error)
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
