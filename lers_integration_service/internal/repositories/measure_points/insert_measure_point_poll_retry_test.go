package measure_points_test

import (
	"context"
	"lers_integration_service/tests"
	"math/rand"
	"testing"
	"time"
)

func TestInsertMeasurePointPollRetry(t *testing.T) {
	ctx := context.Background()
	tests.CleanDb(ctx, appInstance)

	rand.Seed(time.Now().UnixNano())
	originalPollID := rand.Intn(1000) + 1
	retryPollID := rand.Intn(1000) + 1
	status := "PENDING"

	err := insertPollLog(ctx, originalPollID, rand.Intn(1000)+1, rand.Intn(1000)+1, "PENDING")
	if err != nil {
		t.Fatalf("Ошибка вставки оригинального лога poll: %v", err)
	}

	err = repositoryTest.InsertMeasurePointPollRetry(ctx, originalPollID, retryPollID, status)
	if err != nil {
		t.Fatalf("Ошибка вставки retry poll: %v", err)
	}

	validateRetryPollInDB(t, ctx, originalPollID, retryPollID, status)
}

func validateRetryPollInDB(t *testing.T, ctx context.Context, originalPollID int, retryPollID int, expectedStatus string) {
	pgConn := appInstance.S().PostgresDB().DB()

	var (
		dbOriginalPollID int
		dbRetryPollID    int
		dbStatus         string
	)
	err := pgConn.QueryRow(ctx, "SELECT original_poll_id, retrying_poll_id, status FROM measure_points_poll_retry WHERE original_poll_id = $1 AND retrying_poll_id = $2", originalPollID, retryPollID).Scan(&dbOriginalPollID, &dbRetryPollID, &dbStatus)
	if err != nil {
		t.Fatalf("Ошибка запроса к БД: %v", err)
	}

	if dbOriginalPollID != originalPollID || dbRetryPollID != retryPollID || dbStatus != expectedStatus {
		t.Fatalf("Ожидалось original_poll_id=%d, retrying_poll_id=%d, status=%s, но получено original_poll_id=%d, retrying_poll_id=%d, status=%s",
			originalPollID, retryPollID, expectedStatus, dbOriginalPollID, dbRetryPollID, dbStatus)
	}
}
