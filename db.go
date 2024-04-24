package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type EnvDBConfig struct {
	host     string
	port     string
	username string
	password string
	database string
}

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

func ConnectToDB(config EnvDBConfig) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.GetUsername(), config.GetPassword(), config.GetHost(), config.GetPort(), config.GetDatabase())
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type Account struct {
	id      int
	balance float64
}

func CreateAccount(db *sql.DB, id int, balance float64) error {
	query := "INSERT INTO accounts (id, balance) VALUES (?, ?)"
	_, err := db.Exec(query, id, balance)
	if err != nil {
		return err
	}
	return nil
}

func GetAccount(db *sql.DB, id int) (*Account, error) {
	query := "SELECT * FROM accounts WHERE id = ?"
	row := db.QueryRow(query, id)

	account := &Account{}
	err := row.Scan(&account.id, &account.balance)
	if err != nil {
		return nil, err
	}
	return account, nil
}

type Transaction struct {
	TransactionID        string
	SourceAccountID      int
	DestinationAccountID int
	Amount               float64
}

func CreateTransaction(db *sql.DB, transaction *Transaction) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// check that source account id exists and has enough balance
	// Confirm that album inventory is enough for the order.
	var enough bool
	if err = tx.QueryRow("SELECT (balance >= ?) from accounts where id = ?",
		transaction.Amount, transaction.SourceAccountID).Scan(&enough); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("source account does not exist")
		}
		return err
	}
	if !enough {
		return fmt.Errorf("not enough money in balance")
	}

	// check that destination account id exists
	var destinationId int
	err = tx.QueryRow("SELECT id from accounts where id = ?", transaction.DestinationAccountID).Scan(&destinationId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("destination account does not exist")
		}
		return err
	}

	// update both records' balances
	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", transaction.Amount, transaction.SourceAccountID)
	if err != nil {
		return fmt.Errorf("could not update balance for source")
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", transaction.Amount, transaction.DestinationAccountID)
	if err != nil {
		return fmt.Errorf("could not update balance for destination")
	}

	// insert into transactions table
	result, err := tx.Exec("INSERT INTO transactions (source_account_id, destination_account_id, transaction_id, amount) VALUES (?, ?, ?, ?)",
		transaction.SourceAccountID, transaction.DestinationAccountID, transaction.TransactionID, transaction.Amount)
	if err != nil {
		return fmt.Errorf("could not insert transaction")
	}

	// TODO: there's a transaction id here, but need to know what to use it for
	_, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("could not retrieve transaction")
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction")
	}

	return nil
}
