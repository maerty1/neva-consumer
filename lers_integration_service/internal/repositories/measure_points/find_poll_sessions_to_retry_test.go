package measure_points_test

import (
	"context"
	"lers_integration_service/tests"
	"math/rand"
	"testing"
	"time"
)

func TestFindPollSessionsToRetry(t *testing.T) {
	ctx := context.Background()
	tests.CleanDb(ctx, appInstance)

	rand.Seed(time.Now().UnixNano())
	accountID := rand.Intn(1000) + 1
	measurePointID1 := rand.Intn(1000) + 1
	measurePointID2 := rand.Intn(1000) + 1

	err := createAccount(ctx, accountID)
	if err != nil {
		t.Fatalf("Ошибка создания аккаунта: %v", err)
	}

	err = createMeasurePoint(ctx, accountID, measurePointID1)
	if err != nil {
		t.Fatalf("Ошибка создания первой точки учета: %v", err)
	}

	err = createMeasurePoint(ctx, accountID, measurePointID2)
	if err != nil {
		t.Fatalf("Ошибка создания второй точки учета: %v", err)
	}

	err = insertPollSession(ctx, accountID, measurePointID1, "FAILED")
	if err != nil {
		t.Fatalf("Ошибка вставки первой poll сессии: %v", err)
	}

	err = insertPollSession(ctx, accountID, measurePointID2, "SUCCESS")
	if err != nil {
		t.Fatalf("Ошибка вставки второй poll сессии: %v", err)
	}

	pollSessions, err := repositoryTest.FindPollSessionsToRetry(ctx, accountID)
	if err != nil {
		t.Fatalf("Ошибка поиска poll сессий для ретрая: %v", err)
	}

	if len(pollSessions) != 1 {
		t.Fatalf("Ожидалось 1 poll сессия для ретрая, получено: %d", len(pollSessions))
	}

	if pollSessions[0].MeasurePointID != measurePointID1 {
		t.Fatalf("Неправильная measurePointID в результатах: ожидалось %d, получено %d", measurePointID1, pollSessions[0].MeasurePointID)
	}
}

func insertPollSession(ctx context.Context, accountID int, measurePointID int, status string) error {
	pgConn := appInstance.S().PostgresDB().DB()
	_, err := pgConn.Exec(ctx, "INSERT INTO measure_points_poll_log (poll_id, account_id, measure_point_id, status) VALUES ($1, $2, $3, $4)", rand.Intn(1000)+1, accountID, measurePointID, status)
	return err
}
