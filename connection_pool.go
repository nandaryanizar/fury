package fury

import (
	"database/sql"
	"fmt"
)

type connectionPool interface {
	Close() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Ping() error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func newConnectionPool(config *Configuration) (connectionPool, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s ", config.Host, config.Port, config.Username, config.Password, config.DBName)

	if config.SSLMode {
		connStr += "sslmode=enable"
	} else {
		connStr += "sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	return db, nil
}
