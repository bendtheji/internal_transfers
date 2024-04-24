package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"

	apiError "github.com/bendtheji/internal_transfers/errors"
)

type EnvDBConfig struct {
	host     string
	port     string
	username string
	password string
	database string
}

var DbConfig *EnvDBConfig

func NewEnvDBConfig() *EnvDBConfig {
	return &EnvDBConfig{
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		username: os.Getenv("DB_USERNAME"),
		password: os.Getenv("DB_PASSWORD"),
		database: os.Getenv("DB_DATABASE"),
	}
}

func (c *EnvDBConfig) GetHost() string {
	return c.host
}

func (c *EnvDBConfig) GetPort() string {
	return c.port
}

func (c *EnvDBConfig) GetUsername() string {
	return c.username
}

func (c *EnvDBConfig) GetPassword() string {
	return c.password
}

func (c *EnvDBConfig) GetDatabase() string {
	return c.database
}

func InitDbConfig() {
	DbConfig = NewEnvDBConfig()
}

func ConnectToDB(config EnvDBConfig) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.GetUsername(), config.GetPassword(), config.GetHost(), config.GetPort(), config.GetDatabase())
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type Account struct {
	ID      int
	Balance float64
}

func CreateAccount(ctx context.Context, db *sql.DB, id int, balance float64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "INSERT INTO accounts (id, balance) VALUES (?, ?)"
	_, err := db.ExecContext(ctx, query, id, balance)
	if err != nil {
		return apiError.WrapError(err)
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
		return nil, apiError.WrapError(err)
	}
	return account, nil
}

type Transaction struct {
	TransactionID        string
	SourceAccountID      int
	DestinationAccountID int
	Amount               float64
}

func CreateTransaction(ctx context.Context, db *sql.DB, transaction *Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return apiError.WrapError(err)
	}
	defer tx.Rollback()

	// check that source account id exists and has enough balance
	// Confirm that album inventory is enough for the order.
	var enough bool
	if err = tx.QueryRowContext(ctx, "SELECT (balance >= ?) from accounts where id = ? for update",
		transaction.Amount, transaction.SourceAccountID).Scan(&enough); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apiError.WrapError(fmt.Errorf("source account not found: %w", err))
		}
		return apiError.WrapError(err)
	}
	if !enough {
		return apiError.WrapError(apiError.NotEnoughBalanceErr)
	}

	// check that destination account id exists
	var destinationId int
	err = tx.QueryRowContext(ctx, "SELECT id from accounts where id = ? for update", transaction.DestinationAccountID).Scan(&destinationId)
	if err != nil {
		if err == sql.ErrNoRows {
			return apiError.WrapError(fmt.Errorf("destination account not found: %w", err))
		}
		return apiError.WrapError(err)
	}
	// update both records' balances
	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - ? WHERE id = ? ", transaction.Amount, transaction.SourceAccountID)
	if err != nil {
		return apiError.WrapError(fmt.Errorf("could not update balance for source: %w", err))
	}

	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + ? WHERE id = ?", transaction.Amount, transaction.DestinationAccountID)
	if err != nil {
		return apiError.WrapError(fmt.Errorf("could not update balance for destination: %w", err))
	}

	// insert into transactions table
	_, err = tx.ExecContext(ctx, "INSERT INTO transactions (source_account_id, destination_account_id, transaction_id, amount) VALUES (?, ?, ?, ?)",
		transaction.SourceAccountID, transaction.DestinationAccountID, transaction.TransactionID, transaction.Amount)
	if err != nil {
		return apiError.WrapError(fmt.Errorf("could not insert transaction, %w", err))
	}

	if err = tx.Commit(); err != nil {
		return apiError.WrapError(fmt.Errorf("could not commit transaction, %w", err))
	}

	return nil
}
