package db

import (
	"context"
	"database/sql"
	apiError "github.com/bendtheji/internal_transfers/errors"
	"github.com/shopspring/decimal"
	"time"
)

type Account struct {
	ID      int
	Balance decimal.Decimal
}

func CreateAccount(ctx context.Context, db *sql.DB, id int, balance decimal.Decimal) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "INSERT INTO accounts (id, balance) VALUES (?, ?)"
	_, err := db.ExecContext(ctx, query, id, balance)
	if err != nil {
		return apiError.HandleError(err)
	}
	return nil
}

func GetAccount(ctx context.Context, db *sql.DB, id int) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "SELECT * FROM accounts WHERE id = ?"
	row := db.QueryRowContext(ctx, query, id)

	account := &Account{}
	err := row.Scan(&account.ID, &account.Balance)
	if err != nil {
		return nil, apiError.HandleError(err)
	}
	return account, nil
}
