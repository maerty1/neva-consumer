package measure_points_test

import (
	"context"
	"lers_integration_service/internal/repositories/measure_points"
	"lers_integration_service/tests"
	"math/rand"
	"testing"
	"time"
)

func TestUpdateRetryPollStatus(t *testing.T) {
	ctx := context.Background()
	tests.CleanDb(ctx, appInstance)

	rand.Seed(time.Now().UnixNano())
	accountID := rand.Intn(1000) + 1
	measurePointID := rand.Intn(1000) + 1
	originalPollID := rand.Intn(1000) + 1
	retryPollID := rand.Intn(1000) + 1
	initialStatus := "PENDING"
	newStatus := "SUCCESS"

	err := createAccount(ctx, accountID)
	if err != nil {
		t.Fatalf("Ошибка создания аккаунта: %v", err)
	}

	err = createMeasurePoint(ctx, accountID, measurePointID)
	if err != nil {
		t.Fatalf("Ошибка создания точки учета: %v", err)
	}

	err = insertPollLog(ctx, originalPollID, measurePointID, accountID, initialStatus)
	if err != nil {
		t.Fatalf("Ошибка вставки лога poll: %v", err)
	}

	err = insertRetryPoll(ctx, originalPollID, retryPollID, initialStatus)
	if err != nil {
		t.Fatalf("Ошибка вставки retry poll: %v", err)
	}

	retryPoll := measure_points.RetryPollSessions{
		OriginalPollID: originalPollID,
		RetryPollID:    retryPollID,
	}

	err = repositoryTest.UpdateRetryPollStatus(ctx, retryPoll, newStatus)
	if err != nil {
		t.Fatalf("Ошибка обновления статуса retry poll: %v", err)
	}

	validatePollStatusInDB(t, ctx, originalPollID, newStatus)
	validateRetryPollStatusInDB(t, ctx, retryPollID, newStatus)
}

func insertRetryPoll(ctx context.Context, originalPollID int, retryPollID int, status string) error {
	pgConn := appInstance.S().PostgresDB().DB()
	_, err := pgConn.Exec(ctx, "INSERT INTO measure_points_poll_retry (original_poll_id, retrying_poll_id, status) VALUES ($1, $2, $3)", originalPollID, retryPollID, status)
	return err
}

func validateRetryPollStatusInDB(t *testing.T, ctx context.Context, retryPollID int, expectedStatus string) {
	pgConn := appInstance.S().PostgresDB().DB()

	var status string
	err := pgConn.QueryRow(ctx, "SELECT status FROM measure_points_poll_retry WHERE retrying_poll_id = $1", retryPollID).Scan(&status)
	if err != nil {
		t.Fatalf("Ошибка запроса к БД: %v", err)
	}

	if status != expectedStatus {
		t.Fatalf("Ожидалось, что статус будет '%s', но получили '%s'", expectedStatus, status)
	}
}
