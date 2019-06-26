package fury

import (
	"database/sql"
	"fmt"

	// PostgreSQL driver
	_ "github.com/lib/pq"
)

// ConnectionPooler interface
// 	Use this interface as contract for DB connection pool
type ConnectionPooler interface {
	Close() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Ping() error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// NewConnectionPool function
// 	Create connection to DB and return DB connection pool
func NewConnectionPool(config *Configuration) (ConnectionPooler, error) {
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
