package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	apiError "github.com/bendtheji/internal_transfers/errors"
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	TransactionID        string
	SourceAccountID      int
	DestinationAccountID int
	Amount               decimal.Decimal
}

func CreateTransaction(ctx context.Context, db *sql.DB, transaction *Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return apiError.HandleError(err)
	}
	defer tx.Rollback()

	// check that source account id exists and has enough balance
	// we also use the "for update" for the select statement to place a lock
	// on the row that we want to modify
	var enough bool
	if err = tx.QueryRowContext(ctx, "SELECT (balance >= ?) from accounts where id = ? for update",
		transaction.Amount, transaction.SourceAccountID).Scan(&enough); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apiError.HandleError(fmt.Errorf("source account not found: %w", err))
		}
		return apiError.HandleError(err)
	}
	if !enough {
		return apiError.HandleError(apiError.NotEnoughBalanceErr)
	}

	// check that destination account id exists
	// we also use the "for update" for the select statement to place a lock
	// on the row that we want to modify
	var destinationId int
	err = tx.QueryRowContext(ctx, "SELECT id from accounts where id = ? for update", transaction.DestinationAccountID).Scan(&destinationId)
	if err != nil {
		if err == sql.ErrNoRows {
			return apiError.HandleError(fmt.Errorf("destination account not found: %w", err))
		}
		return apiError.HandleError(err)
	}

	// update both records' balances
	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - ? WHERE id = ? ", transaction.Amount, transaction.SourceAccountID)
	if err != nil {
		return apiError.HandleError(fmt.Errorf("could not update balance for source: %w", err))
	}
	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + ? WHERE id = ?", transaction.Amount, transaction.DestinationAccountID)
	if err != nil {
		return apiError.HandleError(fmt.Errorf("could not update balance for destination: %w", err))
	}

	// insert into transactions table
	_, err = tx.ExecContext(ctx, "INSERT INTO transactions (source_account_id, destination_account_id, transaction_id, amount) VALUES (?, ?, ?, ?)",
		transaction.SourceAccountID, transaction.DestinationAccountID, transaction.TransactionID, transaction.Amount)
	if err != nil {
		return apiError.HandleError(fmt.Errorf("could not insert transaction, %w", err))
	}

	// commit
	if err = tx.Commit(); err != nil {
		return apiError.HandleError(fmt.Errorf("could not commit transaction, %w", err))
	}

	return nil
}
