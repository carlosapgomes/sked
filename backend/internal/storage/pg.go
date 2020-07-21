package storage

import (
	"database/sql"

	// blank inport posgres lib
	_ "github.com/lib/pq"
)

// DB type
type DB struct {
	*sql.DB
}

// NewDB initialize a new database connection pool
func NewDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
