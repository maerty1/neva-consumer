package measure_points_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"lers_integration_service/tests"
	"testing"
	"time"
)

func TestInsertMeasurePointData(t *testing.T) {
	ctx := context.Background()
	tests.CleanDb(ctx, appInstance)

	measurePointID := 1
	accountID := 1
	datetime := time.Now().UTC().Format(time.RFC3339)
	values := `{"key": "test values"}`

	err := createAccount(ctx, accountID)
	if err != nil {
		t.Fatalf("Ошибка создания аккаунта: %v", err)
	}

	err = createMeasurePoint(ctx, accountID, measurePointID)
	if err != nil {
		t.Fatalf("Ошибка создания точки учета: %v", err)
	}

	err = repositoryTest.InsertMeasurePointData(ctx, measurePointID, datetime, values)
	if err != nil {
		t.Fatalf("Ошибка вставки данных точки учета: %v", err)
	}

	validateMeasurePointDataInDB(t, ctx, measurePointID, datetime, values)

	updatedValues := `{"key": "updated test values"}`
	err = repositoryTest.InsertMeasurePointData(ctx, measurePointID, datetime, updatedValues)
	if err != nil {
		t.Fatalf("Ошибка обновления данных точки учета: %v", err)
	}

	validateMeasurePointDataInDB(t, ctx, measurePointID, datetime, updatedValues)
}

func validateMeasurePointDataInDB(t *testing.T, ctx context.Context, measurePointID int, datetime string, expectedValues string) {
	pgConn := appInstance.S().PostgresDB().DB()

	query := "SELECT measure_point_id, datetime, values FROM measure_points_data WHERE measure_point_id = $1 AND datetime = $2"
	row := pgConn.QueryRow(ctx, query, measurePointID, datetime)

	var (
		id         int
		dbDatetime time.Time
		values     string
	)

	err := row.Scan(&id, &dbDatetime, &values)
	if err == sql.ErrNoRows {
		t.Fatalf("Не найдены данные для measure_point_id: %d, datetime: %s", measurePointID, datetime)
	} else if err != nil {
		t.Fatalf("Ошибка сканирования строки: %v", err)
	}

	expectedDatetime, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		t.Fatalf("Ошибка парсинга ожидаемого datetime: %v", err)
	}

	if id != measurePointID {
		t.Fatalf("Неверный measure_point_id: ожидается %d, получено %d", measurePointID, id)
	}

	if !dbDatetime.Equal(expectedDatetime) {
		t.Fatalf("Неверный datetime: ожидается %s, получено %s", expectedDatetime.Format(time.RFC3339), dbDatetime.Format(time.RFC3339))
	}

	var expectedMap, valuesMap map[string]interface{}
	if err := json.Unmarshal([]byte(expectedValues), &expectedMap); err != nil {
		t.Fatalf("Ошибка при разборе ожидаемых значений JSON: %v", err)
	}
	if err := json.Unmarshal([]byte(values), &valuesMap); err != nil {
		t.Fatalf("Ошибка при разборе значений JSON из базы данных: %v", err)
	}

	if !equalJSON(expectedMap, valuesMap) {
		t.Fatalf("Неверные values: ожидается %s, получено %s", expectedValues, values)
	}
}

func equalJSON(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}

	return true
}
