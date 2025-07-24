package measure_points_test

import (
	"context"
	"database/sql"
	"lers_integration_service/tests"
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

func TestInsertSyncLog(t *testing.T) {
	ctx := context.Background()
	tests.CleanDb(ctx, appInstance)

	rand.Seed(uint64(time.Now().UnixNano()))
	accountID := rand.Intn(1000) + 1
	measurePointID := rand.Intn(1000) + 1
	logLevel := "INFO"
	logMessage := "test message"

	err := createAccount(ctx, accountID)
	if err != nil {
		t.Fatalf("Ошибка создания аккаунта: %v", err)
	}

	err = createMeasurePoint(ctx, accountID, measurePointID)
	if err != nil {
		t.Fatalf("Ошибка создания точки учета: %v", err)
	}

	err = repositoryTest.InsertSyncLog(ctx, accountID, measurePointID, logLevel, logMessage)
	if err != nil {
		t.Fatalf("Не удалось вставить лог синхронизации: %v", err)
	}

	validateSyncLogInDB(t, ctx, accountID, measurePointID, logLevel, logMessage)
}

func createAccount(ctx context.Context, accountID int) error {
	pgConn := appInstance.S().PostgresDB().DB()
	_, err := pgConn.Exec(ctx, "INSERT INTO accounts (id, name, token, server_host) VALUES ($1, 'test', 'test', 'test')", accountID)
	return err
}

func createMeasurePoint(ctx context.Context, accountID int, measurePointID int) error {
	pgConn := appInstance.S().PostgresDB().DB()
	_, err := pgConn.Exec(ctx, "INSERT INTO measure_points (id, account_id, title) VALUES ($1, $2, 'test')", measurePointID, accountID)
	return err
}

func validateSyncLogInDB(t *testing.T, ctx context.Context, accountID int, measurePointID int, logLevel string, logMessage string) {
	pgConn := appInstance.S().PostgresDB().DB()

	rows, err := pgConn.Query(ctx, "SELECT account_id, measure_point_id, level, message FROM accounts_sync_log")
	if err != nil {
		t.Fatalf("Ошибка запроса к БД: %v", err)
	}
	defer rows.Close()

	var count int
	var (
		account_id       int
		measure_point_id sql.NullInt32
		level            string
		message          string
		found            bool
	)

	for rows.Next() {
		count++
		err = rows.Scan(&account_id, &measure_point_id, &level, &message)
		if err != nil {
			t.Fatalf("Ошибка сканирования строк: %v", err)
		}

		if account_id == accountID && measure_point_id.Int32 == int32(measurePointID) && level == logLevel && message == logMessage {
			found = true
			break
		}
	}

	if err := rows.Err(); err != nil {
		t.Fatalf("Ошибка итерации по строкам: %v", err)
	}

	if count == 0 {
		t.Fatalf("Неверное количество записей в БД: %d", count)
	}

	if !found {
		t.Fatalf("Не удалось найти запись с нужными данными")
	}
}
