package zulu

import (
	"context"
	"time"
	"zulu_updater/internal/db"
	"zulu_updater/internal/models"
)

type Repository interface {
	SelectZwsTypeByElemId(ctx context.Context, elemId int) (int, error)
	SelectValNameByZwsType(ctx context.Context, zwsType int, extractionType string) ([]string, error)
	InsertRecords(ctx context.Context, fields []models.Field, elemId int) error
	InsertObjectRecordsJson(ctx context.Context, data models.Records, updateType string, elemID int, lersTS time.Time) (*ObjectRecordsFromJsonResponse, error)
}

type repository struct {
	db db.PostgresClient
}

func NewRepository(db db.PostgresClient) *repository {
	return &repository{
		db: db,
	}
}
