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
	ID        int64       `json:"id" db:"id"`
	UserID    int64       `json:"user_id" db:"user_id"`
	Status    OrderStatus `json:"status" db:"status"`
	Total     float64     `json:"total" db:"total"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	Items     []OrderItem `json:"items"`
}

type OrderItem struct {
	ID        int64   `json:"id" db:"id"`
	OrderID   int64   `json:"order_id" db:"order_id"`
	ProductID int64   `json:"product_id" db:"product_id"`
	Quantity  int     `json:"quantity" db:"quantity"`
	Price     float64 `json:"price" db:"price"`
}

type CartItem struct {
	ID        int64 `json:"id" db:"id"`
	UserID    int64 `json:"user_id" db:"user_id"`
	ProductID int64 `json:"product_id" db:"product_id"`
	Quantity  int   `json:"quantity" db:"quantity"`
}
