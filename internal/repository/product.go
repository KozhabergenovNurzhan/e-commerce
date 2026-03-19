package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"ecommerce/internal/models"
)

type ProductRepository interface {
	Create(ctx context.Context, p *models.Product) error
	GetByID(ctx context.Context, id int64) (*models.Product, error)
	List(ctx context.Context, f ProductFilter) ([]models.Product, int, error)
	Update(ctx context.Context, p *models.Product) error
	Delete(ctx context.Context, id int64) error
	DecreaseStock(ctx context.Context, id int64, qty int) error
}

type ProductFilter struct {
	Search     string
	CategoryID int64
	Limit      int
	Offset     int
}

type ProductRepo struct {
	db *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(ctx context.Context, p *models.Product) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO products (name, description, price, stock, category_id)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`,
		p.Name, p.Description, p.Price, p.Stock, p.CategoryID,
	).Scan(&p.ID, &p.CreatedAt)
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

func (r *ProductRepo) List(ctx context.Context, f ProductFilter) ([]models.Product, int, error) {
	args := []any{}
	where := []string{}
	i := 1

	if f.Search != "" {
		where = append(where, fmt.Sprintf("name ILIKE $%d", i))
		args = append(args, "%"+f.Search+"%")
		i++
	}

	if f.CategoryID > 0 {
		where = append(where, fmt.Sprintf("category_id = $%d", i))
		args = append(args, f.CategoryID)
		i++
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}

	query := fmt.Sprintf(
		`SELECT * FROM products %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		whereClause, i, i+1,
	)
	args = append(args, f.Limit, f.Offset)

	var products []models.Product
	if err := r.db.SelectContext(ctx, &products, query, args...); err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}

	return products, total, nil
}

func (r *ProductRepo) Update(ctx context.Context, p *models.Product) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE products SET name=$1, description=$2, price=$3, stock=$4, category_id=$5 WHERE id=$6`,
		p.Name, p.Description, p.Price, p.Stock, p.CategoryID, p.ID,
	)
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
