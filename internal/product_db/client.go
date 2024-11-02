package product_db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Database struct to manage database connection pool
type Database struct {
	pool *pgxpool.Pool
}

// NewDatabase initializes the database connection pool
func NewDatabase(dsn string) (*Database, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return &Database{pool: pool}, nil
}

// Close closes the database connection pool
func (db *Database) Close() {
	db.pool.Close()
}
