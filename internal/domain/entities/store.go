package entities

import "time"

type Store struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     uint64    `json:"owner_id"` // ID владельца магазина
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
