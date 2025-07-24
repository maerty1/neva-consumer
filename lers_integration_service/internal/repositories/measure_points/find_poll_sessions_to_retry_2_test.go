package measure_points_test

import (
	"context"
	"fmt"
	"lers_integration_service/tests"
	"math/rand"
	"testing"
	"time"
)

func TestFindPollSessionsToRetry2(t *testing.T) {
	ctx := context.Background()
	tests.CleanDb(ctx, appInstance)

	rand.Seed(time.Now().UnixNano())
	accountID := rand.Intn(1000) + 1
	measurePointID := rand.Intn(1000) + 1
	pollID1 := rand.Intn(1000) + 1
	pollID2 := pollID1 + rand.Intn(1000) + 1
	pollID3 := pollID2 + rand.Intn(1000) + 1

	err := createAccount(ctx, accountID)
	if err != nil {
		t.Fatalf("Ошибка создания аккаунта: %v", err)
	}

	err = createMeasurePoint(ctx, accountID, measurePointID)
	if err != nil {
		t.Fatalf("Ошибка создания точки учета: %v", err)
	}

	err = insertPollLog(ctx, pollID1, measurePointID, accountID, "PENDING")
	if err != nil {
		t.Fatalf("Ошибка вставки первой poll сессии: %v", err)
	}

	err = insertPollLog(ctx, pollID2, measurePointID, accountID, "PENDING")
	if err != nil {
		t.Fatalf("Ошибка вставки второй poll сессии: %v", err)
	}

	err = insertPollLog(ctx, pollID3, measurePointID, accountID, "FAILED")
	if err != nil {
		t.Fatalf("Ошибка вставки третьей poll сессии: %v", err)
	}

	pollSessions, err := repositoryTest.FindPollSessionsToRetry2(ctx, accountID)
	if err != nil {
		t.Fatalf("Ошибка поиска poll сессий для ретрая: %v", err)
	}

	expectedPollCount := 1
	if len(pollSessions) != expectedPollCount {
		t.Fatalf("Ожидалось %d poll сессий для ретрая, получено: %d", expectedPollCount, len(pollSessions))
	}

	fmt.Println("---------------")
	fmt.Println(pollID1)
	fmt.Println(pollID2)
	fmt.Println(pollID3)
	fmt.Println("---------------")

	if pollSessions[0].PollID != pollID3 || pollSessions[0].MeasurePointID != measurePointID {
		t.Fatalf("Неверная poll сессия в результатах: ожидалось pollID %d и measurePointID %d, получено pollID %d и measurePointID %d",
			pollID3, measurePointID, pollSessions[0].PollID, pollSessions[0].MeasurePointID)
	}

	// Проверяем, что предыдущие записи были помечены как "OUTDATED"
	verifyOutdatedStatus(t, ctx, pollID1)
	verifyOutdatedStatus(t, ctx, pollID2)
}

func verifyOutdatedStatus(t *testing.T, ctx context.Context, pollID int) {
	pgConn := appInstance.S().PostgresDB().DB()
	var status string
	err := pgConn.QueryRow(ctx, "SELECT status FROM measure_points_poll_log WHERE poll_id = $1", pollID).Scan(&status)
	if err != nil {
		t.Fatalf("Ошибка запроса к БД: %v", err)
	}
	if status != "OUTDATED" {
		t.Fatalf("Ожидалось, что статус будет 'OUTDATED' для pollID %d, но получен статус '%s'", pollID, status)
	}
}
