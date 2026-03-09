package models

import "time"

type Product struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	Stock       int       `json:"stock" db:"stock"`
	CategoryID  int64     `json:"category_id" db:"category_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Category struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
