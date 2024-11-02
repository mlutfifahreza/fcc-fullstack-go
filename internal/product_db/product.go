package product_db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type Product struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

func (db *Database) CreateProduct(ctx context.Context, product *Product) error {
	query := `INSERT INTO products (id, name) VALUES ($1, $2)`
	_, err := db.pool.Exec(ctx, query, product.ID, product.Name)
	if err != nil {
		return fmt.Errorf("error creating product: %v", err)
	}
	return nil
}

func (db *Database) GetProduct(ctx context.Context, id string) (*Product, error) {
	query := `SELECT id, name FROM products WHERE id = $1`
	row := db.pool.QueryRow(ctx, query, id)

	var product Product
	if err := row.Scan(&product.ID, &product.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("error getting product: %v", err)
	}
	return &product, nil
}

func (db *Database) UpdateProduct(ctx context.Context, product *Product) error {
	query := `UPDATE products SET name = $1 WHERE id = $2`
	result, err := db.pool.Exec(ctx, query, product.Name, product.ID)
	if err != nil {
		return fmt.Errorf("error updating product: %v", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) DeleteProduct(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := db.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting product: %v", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
