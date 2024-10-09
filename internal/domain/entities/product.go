package entities

import "time"

type Product struct {
	ID          uint64
	Name        string
	Description string
	Price       float64
	OwnerID     uint64
	StoreID     uint64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      string // Доступен, Продан и т.д.
}
