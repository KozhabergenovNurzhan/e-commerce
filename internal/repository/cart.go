package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"ecommerce/internal/models"
)

type CartRepo struct {
	db *sqlx.DB
}

func NewCartRepo(db *sqlx.DB) *CartRepo {
	return &CartRepo{db: db}
}

func (r *CartRepo) GetByUser(ctx context.Context, userID int64) ([]models.CartItem, error) {
	var items []models.CartItem
	err := r.db.SelectContext(ctx, &items, `SELECT * FROM cart_items WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("get cart: %w", err)
	}
	return items, nil
}

func (r *CartRepo) Upsert(ctx context.Context, item *models.CartItem) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO cart_items (user_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity`,
		item.UserID, item.ProductID, item.Quantity,
	)
	return err
}

func (r *CartRepo) Remove(ctx context.Context, userID, productID int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM cart_items WHERE user_id = $1 AND product_id = $2`,
		userID, productID,
	)
	return err
}

func (r *CartRepo) Clear(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM cart_items WHERE user_id = $1`, userID)
	return err
}
