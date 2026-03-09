package models

import "time"

type Product struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float64   `db:"price"`
	Stock       int       `db:"stock"`
	CategoryID  int64     `db:"category_id"`
	CreatedAt   time.Time `db:"created_at"`
}

type Category struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
