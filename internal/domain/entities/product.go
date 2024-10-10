package entities

import "time"

type Product struct {
	ID          uint64    `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price,omitempty"`
	OwnerID     uint64    `json:"owner_id,omitempty"`
	StoreID     uint64    `json:"store_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      string    `json:"status,omitempty"` // Доступен, Продан и т.д.
}
