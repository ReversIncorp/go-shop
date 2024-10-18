package entities

import "time"

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Price       float64   `json:"price" validate:"required"`
	Quantity    int       `json:"quantity" validate:"required"`
	CategoryID  int64     `json:"category_id" validate:"required"`
	StoreID     int64     `json:"store_id" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
