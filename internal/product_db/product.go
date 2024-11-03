package product_db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

type Product struct {
	ID        string    `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type GetProductListFilter struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
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
	query := `SELECT id, name, created_at, updated_at FROM products WHERE id = $1`
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

func (db *Database) GetProductList(ctx context.Context, filter GetProductListFilter) ([]Product, error) {
	query := `SELECT id, name, created_at, updated_at FROM products ORDER BY created_at ASC LIMIT $1 OFFSET $2`
	rows, err := db.pool.Query(ctx, query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, fmt.Errorf("error getting product list: %v", err)
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning product: %v", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over product rows: %v", err)
	}

	return products, nil
}

func (db *Database) UpdateProduct(ctx context.Context, product *Product) error {
	query := `UPDATE products SET name = $1, updated_at = $2 WHERE id = $3`
	result, err := db.pool.Exec(ctx, query, product.Name, product.UpdatedAt, product.ID)
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
