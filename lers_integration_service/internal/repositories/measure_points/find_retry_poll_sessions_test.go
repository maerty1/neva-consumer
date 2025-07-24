package measure_points_test

import (
	"context"
	"lers_integration_service/tests"
	"math/rand"
	"testing"
	"time"
)

func TestFindRetryPollSessions(t *testing.T) {
	ctx := context.Background()
	tests.CleanDb(ctx, appInstance)

	rand.Seed(time.Now().UnixNano())
	accountID := rand.Intn(1000) + 1
	originalPollID1 := rand.Intn(1000) + 1
	retryPollID1 := rand.Intn(1000) + 1
	originalPollID2 := rand.Intn(1000) + 1
	retryPollID2 := rand.Intn(1000) + 1

	err := createAccount(ctx, accountID)
	if err != nil {
		t.Fatalf("Ошибка создания аккаунта: %v", err)
	}

	err = createPollLog(ctx, accountID, originalPollID1, "IN_PROGRESS")
	if err != nil {
		t.Fatalf("Ошибка создания первой poll записи: %v", err)
	}

	err = createPollRetry(ctx, originalPollID1, retryPollID1, "IN_PROGRESS")
	if err != nil {
		t.Fatalf("Ошибка создания первой retry poll записи: %v", err)
	}

	err = createPollLog(ctx, accountID, originalPollID2, "IN_PROGRESS")
	if err != nil {
		t.Fatalf("Ошибка создания второй poll записи: %v", err)
	}

	err = createPollRetry(ctx, originalPollID2, retryPollID2, "SUCCESS")
	if err != nil {
		t.Fatalf("Ошибка создания второй retry poll записи: %v", err)
	}

	err = createPollRetry(ctx, originalPollID2, retryPollID2, "FAILED")
	if err != nil {
		t.Fatalf("Ошибка создания третьей retry poll записи: %v", err)
	}

	pollSessions, err := repositoryTest.FindRetryPollSessions(ctx, accountID)
	if err != nil {
		t.Fatalf("Ошибка поиска retry poll сессий: %v", err)
	}

	if len(pollSessions) != 1 {
		t.Fatalf("Ожидалось 1 retry poll сессия для ретрая, получено: %d", len(pollSessions))
	}

	if pollSessions[0].OriginalPollID != originalPollID1 || pollSessions[0].RetryPollID != retryPollID1 {
		t.Fatalf("Неверная retry poll сессия в результатах: ожидалось оригинальный poll ID %d и retry poll ID %d, получено оригинальный poll ID %d и retry poll ID %d",
			originalPollID1, retryPollID1, pollSessions[0].OriginalPollID, pollSessions[0].RetryPollID)
	}
}

func createPollLog(ctx context.Context, accountID int, pollID int, status string) error {
	pgConn := appInstance.S().PostgresDB().DB()
	_, err := pgConn.Exec(ctx, "INSERT INTO measure_points_poll_log (poll_id, account_id, measure_point_id, status) VALUES ($1, $2, $3, $4)", pollID, accountID, rand.Intn(1000)+1, status)
	return err
}

func createPollRetry(ctx context.Context, originalPollID int, retryPollID int, status string) error {
	pgConn := appInstance.S().PostgresDB().DB()
	_, err := pgConn.Exec(ctx, "INSERT INTO measure_points_poll_retry (original_poll_id, retrying_poll_id, status) VALUES ($1, $2, $3)", originalPollID, retryPollID, status)
	return err
}
