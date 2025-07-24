package measure_points_test

import (
	"context"
	"lers_integration_service/tests"
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

func TestInsertMeasurePoint(t *testing.T) {
	ctx := context.Background()
	tests.CleanDb(ctx, appInstance)

	rand.Seed(uint64(time.Now().UnixNano()))
	accountID := rand.Intn(1000) + 1
	measurePointID := rand.Intn(1000) + 1
	deviceID := rand.Intn(1000) + 1

	initialTitle := "Initial Title"
	updatedTitle := "Updated Title"
	initialAddress := "Inital Address"
	updatedAddress := "Updated Address"
	initialFullTitle := "Initial Full Title"
	updatedFullTitle := "Updated Full Title"

	err := createAccount(ctx, accountID)
	if err != nil {
		t.Fatalf("Ошибка создания аккаунта: %v", err)
	}

	err = repositoryTest.InsertMeasurePoint(ctx, accountID, measurePointID, deviceID, initialTitle, initialFullTitle, initialAddress, "")
	if err != nil {
		t.Fatalf("Ошибка вставки точки учета: %v", err)
	}

	validateMeasurePointInDB(t, ctx, measurePointID, accountID, initialTitle)

	err = repositoryTest.InsertMeasurePoint(ctx, accountID, measurePointID, deviceID, updatedTitle, updatedFullTitle, updatedAddress, "")
	if err != nil {
		t.Fatalf("Ошибка обновления точки учета: %v", err)
	}

	validateMeasurePointInDB(t, ctx, measurePointID, accountID, updatedTitle)
}

func validateMeasurePointInDB(t *testing.T, ctx context.Context, measurePointID int, accountID int, expectedTitle string) {
	pgConn := appInstance.S().PostgresDB().DB()

	query := "SELECT id, account_id, title FROM measure_points WHERE id = $1"
	row := pgConn.QueryRow(ctx, query, measurePointID)

	var (
		id          int
		dbAccountID int
		title       string
	)

	err := row.Scan(&id, &dbAccountID, &title)
	if err != nil {
		t.Fatalf("Ошибка сканирования строки: %v", err)
	}

	if id != measurePointID {
		t.Fatalf("Неверный id: ожидается %d, получено %d", measurePointID, id)
	}

	if dbAccountID != accountID {
		t.Fatalf("Неверный account_id: ожидается %d, получено %d", accountID, dbAccountID)
	}

	if title != expectedTitle {
		t.Fatalf("Неверный title: ожидается %s, получено %s", expectedTitle, title)
	}
}
