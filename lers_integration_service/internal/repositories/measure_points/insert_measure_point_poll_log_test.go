package measure_points_test

import (
	"context"
	"lers_integration_service/tests"
	"math/rand"
	"testing"
	"time"
)

func TestInsertMeasurePointPollLog(t *testing.T) {
	ctx := context.Background()
	tests.CleanDb(ctx, appInstance)

	rand.Seed(time.Now().UnixNano())
	accountID := rand.Intn(1000) + 1
	measurePointID := rand.Intn(1000) + 1
	pollID := rand.Intn(1000) + 1
	message := "Test log message"

	err := createAccount(ctx, accountID)
	if err != nil {
		t.Fatalf("Ошибка создания аккаунта: %v", err)
	}

	err = createMeasurePoint(ctx, accountID, measurePointID)
	if err != nil {
		t.Fatalf("Ошибка создания точки учета: %v", err)
	}

	err = repositoryTest.InsertMeasurePointPollLog(ctx, pollID, measurePointID, accountID, message)
	if err != nil {
		t.Fatalf("Ошибка вставки лога точки учета: %v", err)
	}

	validateMeasurePointPollLogInDB(t, ctx, pollID, measurePointID, accountID, message)
}

func validateMeasurePointPollLogInDB(t *testing.T, ctx context.Context, pollID int, measurePointID int, accountID int, message string) {
	pgConn := appInstance.S().PostgresDB().DB()

	rows, err := pgConn.Query(ctx, "SELECT poll_id, message, measure_point_id, account_id, status FROM measure_points_poll_log WHERE poll_id = $1", pollID)
	if err != nil {
		t.Fatalf("Ошибка запроса к БД: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Fatalf("Не найдено записей с poll_id: %d", pollID)
	}

	var (
		dbPollID         int
		dbMessage        string
		dbMeasurePointID int
		dbAccountID      int
		dbStatus         string
	)

	err = rows.Scan(&dbPollID, &dbMessage, &dbMeasurePointID, &dbAccountID, &dbStatus)
	if err != nil {
		t.Fatalf("Ошибка сканирования строки: %v", err)
	}

	if dbPollID != pollID || dbMessage != message || dbMeasurePointID != measurePointID || dbAccountID != accountID || dbStatus != "SUCCESS" {
		t.Fatalf("Данные в БД не соответствуют ожиданиям. Получено: poll_id=%d, message=%s, measure_point_id=%d, account_id=%d, status=%s",
			dbPollID, dbMessage, dbMeasurePointID, dbAccountID, dbStatus)
	}

	if rows.Next() {
		t.Fatalf("Найдено больше одной записи с poll_id: %d", pollID)
	}
}
