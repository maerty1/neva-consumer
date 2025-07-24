package message

import (
	"context"
	"scada_consumer/internal/db"
	"scada_consumer/internal/message_broker/models"
)

type Repository interface {
	WriteBatch(ctx context.Context, messages []models.Message) error
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
