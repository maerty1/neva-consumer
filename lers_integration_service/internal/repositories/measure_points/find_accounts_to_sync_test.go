package measure_points_test

import (
	"context"
	"lers_integration_service/internal/models"
	"lers_integration_service/tests"
	"testing"
)

func TestFindAccountsToSync(t *testing.T) {
	ctx := context.Background()
	tests.CleanDb(ctx, appInstance)

	accounts := []models.AccountToSync{
		{
			ID:         1,
			Token:      "token1",
			ServerHost: "server1",
		},
		{
			ID:         2,
			Token:      "token2",
			ServerHost: "server2",
		},
		{
			ID:         3,
			Token:      "token3",
			ServerHost: "server3",
		},
	}

	for _, account := range accounts {
		err := createAccountWithSyncDatetime(ctx, account.ID, account.Token, account.ServerHost)
		if err != nil {
			t.Fatalf("Ошибка создания аккаунта: %v", err)
		}
	}

	foundAccounts, err := repositoryTest.FindAccountsToSync(ctx)
	if err != nil {
		t.Fatalf("Ошибка вызова FindAccountsToSync: %v", err)
	}

	if len(foundAccounts) != len(accounts) {
		t.Fatalf("Неверное количество найденных аккаунтов: ожидается %d, получено %d", len(accounts), len(foundAccounts))
	}

	for i, account := range accounts {
		if foundAccounts[i].ID != account.ID {
			t.Fatalf("Неверный ID аккаунта: ожидается %d, получено %d", account.ID, foundAccounts[i].ID)
		}
		if foundAccounts[i].Token != account.Token {
			t.Fatalf("Неверный Token аккаунта: ожидается %s, получено %s", account.Token, foundAccounts[i].Token)
		}
		if foundAccounts[i].ServerHost != account.ServerHost {
			t.Fatalf("Неверный ServerHost аккаунта: ожидается %s, получено %s", account.ServerHost, foundAccounts[i].ServerHost)
		}
	}
}

func createAccountWithSyncDatetime(ctx context.Context, accountID int, token string, serverHost string) error {
	pgConn := appInstance.S().PostgresDB().DB()
	_, err := pgConn.Exec(ctx, "INSERT INTO accounts (id, token, name, server_host) VALUES ($1, $2, $3, $4)", accountID, token, "test", serverHost)
	return err
}
