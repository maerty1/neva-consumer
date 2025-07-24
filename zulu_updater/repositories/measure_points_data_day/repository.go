package measure_points_data_day

import (
	"context"
	"time"
	"zulu_updater/internal/db"
)

type Repository interface {
	GetDataParameterByDay(ctx context.Context, day time.Time, paramName string) (float64, error)
}

type repository struct {
	db db.PostgresClient
}

func NewRepository(db db.PostgresClient) *repository {
	return &repository{
		db: db,
	}
}
