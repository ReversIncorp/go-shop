package entities

import "time"

type Product struct {
	ID          uint64    `gorm:"primaryKey"`
	Name        string    `gorm:"size:100;not null" json:"name" validate:"required"`
	Description string    `gorm:"type:text" json:"description" validate:"required"`
	Price       float64   `gorm:"not null" json:"price" validate:"required"`
	Quantity    int       `gorm:"not null" json:"quantity" validate:"required"`
	CategoryID  uint64    `gorm:"not null" json:"category_id" validate:"required"` // Внешний ключ для категории
	Category    Category  `gorm:"foreignKey:CategoryID" json:"category"`           // Связь с категорией
	StoreID     uint64    `gorm:"not null" json:"store_id" validate:"required"`    // Внешний ключ для магазина
	Store       Store     `gorm:"foreignKey:StoreID" json:"store"`                 // Связь с магазином
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
