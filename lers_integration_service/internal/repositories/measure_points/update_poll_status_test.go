package measure_points_test

import (
	"context"
	"lers_integration_service/tests"
	"math/rand"
	"testing"
	"time"
)

func TestUpdatePollStatus(t *testing.T) {
	ctx := context.Background()
	tests.CleanDb(ctx, appInstance)

	rand.Seed(time.Now().UnixNano())
	accountID := rand.Intn(1000) + 1
	measurePointID := rand.Intn(1000) + 1
	pollID := rand.Intn(1000) + 1
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

	err = insertPollLog(ctx, pollID, measurePointID, accountID, initialStatus)
	if err != nil {
		t.Fatalf("Ошибка вставки лога poll: %v", err)
	}

	err = repositoryTest.UpdatePollStatus(ctx, pollID, newStatus)
	if err != nil {
		t.Fatalf("Ошибка обновления статуса poll: %v", err)
	}

	validatePollStatusInDB(t, ctx, pollID, newStatus)
}

func insertPollLog(ctx context.Context, pollID int, measurePointID int, accountID int, status string) error {
	pgConn := appInstance.S().PostgresDB().DB()
	_, err := pgConn.Exec(ctx, "INSERT INTO measure_points_poll_log (poll_id, measure_point_id, account_id, status) VALUES ($1, $2, $3, $4)", pollID, measurePointID, accountID, status)
	return err
}

func validatePollStatusInDB(t *testing.T, ctx context.Context, pollID int, expectedStatus string) {
	pgConn := appInstance.S().PostgresDB().DB()

	var status string
	err := pgConn.QueryRow(ctx, "SELECT status FROM measure_points_poll_log WHERE poll_id = $1", pollID).Scan(&status)
	if err != nil {
		t.Fatalf("Ошибка запроса к БД: %v", err)
	}

	if status != expectedStatus {
		t.Fatalf("Ожидалось, что статус будет '%s', но получили '%s'", expectedStatus, status)
	}
}
