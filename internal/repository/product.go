package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"ecommerce/internal/models"
)

type ProductRepo struct {
	db *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(ctx context.Context, p *models.Product) error {
	query := `
		INSERT INTO products (name, description, price, stock, category_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	return r.db.QueryRowContext(ctx, query, p.Name, p.Description, p.Price, p.Stock, p.CategoryID).
		Scan(&p.ID, &p.CreatedAt)
}

func (r *ProductRepo) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	var p models.Product
	err := r.db.GetContext(ctx, &p, `SELECT * FROM products WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}
	return &p, nil
}

func (r *ProductRepo) GetAll(ctx context.Context, limit, offset int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.SelectContext(ctx, &products,
		`SELECT * FROM products ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	return products, nil
}

func (r *ProductRepo) Update(ctx context.Context, p *models.Product) error {
	query := `
		UPDATE products SET name=$1, description=$2, price=$3, stock=$4, category_id=$5
		WHERE id=$6`
	_, err := r.db.ExecContext(ctx, query, p.Name, p.Description, p.Price, p.Stock, p.CategoryID, p.ID)
	return err
}

func (r *ProductRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
	return err
}

func (r *ProductRepo) DecreaseStock(ctx context.Context, id int64, qty int) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1`,
		qty, id,
	)
	if err != nil {
		return fmt.Errorf("decrease stock: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return models.ErrInsufficientStock
	}
	return nil
}
