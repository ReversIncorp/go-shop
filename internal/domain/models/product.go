package models

import "time"

type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	StoreID     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      string // Доступен, Продан и т.д.
}
