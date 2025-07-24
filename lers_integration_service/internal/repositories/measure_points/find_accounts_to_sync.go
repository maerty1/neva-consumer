package measure_points

import (
	"context"
	"lers_integration_service/internal/models"
)

func (r *repository) FindAccountsToSync(ctx context.Context) ([]models.AccountToSync, error) {
	query := `
		SELECT id, token, server_host
		FROM accounts
	`

	rows, err := r.db.DB().Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.AccountToSync
	for rows.Next() {
		var account models.AccountToSync
		err = rows.Scan(&account.ID, &account.Token, &account.ServerHost)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}
