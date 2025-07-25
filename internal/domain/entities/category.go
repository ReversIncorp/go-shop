package entities

import "time"

type Category struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	Stores    []Store   `gorm:"many2many:store_categories" json:"stores,omitempty"`
	Products  []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}
