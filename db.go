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
