package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"ecommerce/internal/models"
)

type CategoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepo(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) Create(ctx context.Context, c *models.Category) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO categories (name) VALUES ($1) RETURNING id`,
		c.Name,
	).Scan(&c.ID)
}

func (r *CategoryRepo) GetAll(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.SelectContext(ctx, &categories, `SELECT * FROM categories ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	return categories, nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id int64) (*models.Category, error) {
	var c models.Category
	err := r.db.GetContext(ctx, &c, `SELECT * FROM categories WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get category: %w", err)
	}
	return &c, nil
}

func (r *CategoryRepo) Update(ctx context.Context, c *models.Category) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE categories SET name = $1 WHERE id = $2`,
		c.Name, c.ID,
	)
	return err
}

func (r *CategoryRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE id = $1`, id)
	return err
}
