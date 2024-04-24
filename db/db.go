package db

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

func (c *EnvDBConfig) getHost() string {
	return c.host
}

func (c *EnvDBConfig) getPort() string {
	return c.port
}

func (c *EnvDBConfig) getUsername() string {
	return c.username
}

func (c *EnvDBConfig) getPassword() string {
	return c.password
}

func (c *EnvDBConfig) getDatabase() string {
	return c.database
}

func InitDbConfig() {
	DbConfig = NewEnvDBConfig()
}

func ConnectToDB(config EnvDBConfig) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.getUsername(), config.getPassword(), config.getHost(), config.getPort(), config.getDatabase())
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
