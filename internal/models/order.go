package models

import "time"

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPaid      OrderStatus = "paid"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
)

type Order struct {
	ID        int64       `db:"id"`
	UserID    int64       `db:"user_id"`
	Status    OrderStatus `db:"status"`
	Total     float64     `db:"total"`
	CreatedAt time.Time   `db:"created_at"`
	Items     []OrderItem
}

type OrderItem struct {
	ID        int64   `db:"id"`
	OrderID   int64   `db:"order_id"`
	ProductID int64   `db:"product_id"`
	Quantity  int     `db:"quantity"`
	Price     float64 `db:"price"`
}

type CartItem struct {
	ID        int64 `db:"id"`
	UserID    int64 `db:"user_id"`
	ProductID int64 `db:"product_id"`
	Quantity  int   `db:"quantity"`
}
