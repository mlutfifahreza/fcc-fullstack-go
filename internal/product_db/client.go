package product_db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func NewDatabase(dsn string) (*Database, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return &Database{pool: pool}, nil
}

func (db *Database) Close() {
	db.pool.Close()
}
