package measure_points_test

import (
	"context"
	"lers_integration_service/tests"
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

func TestGetLastMeasurePointDatetime(t *testing.T) {
	ctx := context.Background()
	tests.CleanDb(ctx, appInstance)

	rand.Seed(uint64(time.Now().UnixNano()))
	accountID := rand.Intn(1000) + 1
	measurePointID := rand.Intn(1000) + 1

	err := createAccount(ctx, accountID)
	if err != nil {
		t.Fatalf("Ошибка создания аккаунта: %v", err)
	}

	err = createMeasurePoint(ctx, accountID, measurePointID)
	if err != nil {
		t.Fatalf("Ошибка создания точки учета: %v", err)
	}

	datetime1 := time.Now().UTC().Add(-2 * time.Hour).Format(time.RFC3339)
	values1 := `{"key": "value1"}`
	err = repositoryTest.InsertMeasurePointData(ctx, measurePointID, datetime1, values1)
	if err != nil {
		t.Fatalf("Ошибка вставки данных точки учета: %v", err)
	}

	datetime2 := time.Now().UTC().Add(-1 * time.Hour).Format(time.RFC3339)
	values2 := `{"key": "value2"}`
	err = repositoryTest.InsertMeasurePointData(ctx, measurePointID, datetime2, values2)
	if err != nil {
		t.Fatalf("Ошибка вставки данных точки учета: %v", err)
	}

	datetime3 := time.Now().UTC().Format(time.RFC3339)
	values3 := `{"key": "value3"}`
	err = repositoryTest.InsertMeasurePointData(ctx, measurePointID, datetime3, values3)
	if err != nil {
		t.Fatalf("Ошибка вставки данных точки учета: %v", err)
	}

	lastDatetime, err := repositoryTest.GetLastMeasurePointDatetime(ctx, measurePointID)
	if err != nil {
		t.Fatalf("Ошибка получения последнего datetime: %v", err)
	}

	if lastDatetime != datetime3 {
		t.Fatalf("Неверный последний datetime: ожидается %s, получено %s", datetime3, lastDatetime)
	}
}
