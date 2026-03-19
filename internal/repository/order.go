package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"ecommerce/internal/models"
)

type OrderRepository interface {
	Create(ctx context.Context, o *models.Order) error
	GetByID(ctx context.Context, id int64) (*models.Order, error)
	ListByUser(ctx context.Context, userID int64) ([]models.Order, error)
}

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(ctx context.Context, o *models.Order) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx,
		`INSERT INTO orders (user_id, status, total) VALUES ($1, $2, $3) RETURNING id, created_at`,
		o.UserID, o.Status, o.Total,
	).Scan(&o.ID, &o.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	for i := range o.Items {
		o.Items[i].OrderID = o.ID
		err = tx.QueryRowContext(ctx,
			`INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4) RETURNING id`,
			o.Items[i].OrderID, o.Items[i].ProductID, o.Items[i].Quantity, o.Items[i].Price,
		).Scan(&o.Items[i].ID)
		if err != nil {
			return fmt.Errorf("insert order item: %w", err)
		}
	}

	return tx.Commit()
}

func (r *OrderRepo) GetByID(ctx context.Context, id int64) (*models.Order, error) {
	var o models.Order
	err := r.db.GetContext(ctx, &o, `SELECT * FROM orders WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	err = r.db.SelectContext(ctx, &o.Items, `SELECT * FROM order_items WHERE order_id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("get order items: %w", err)
	}

	return &o, nil
}

func (r *OrderRepo) ListByUser(ctx context.Context, userID int64) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.SelectContext(ctx, &orders,
		`SELECT * FROM orders WHERE user_id = $1 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	return orders, nil
}
