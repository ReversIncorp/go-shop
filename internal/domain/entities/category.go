package entities

import "time"

type Category struct {
	ID          uint64    `gorm:"primaryKey"`
	Name        string    `gorm:"size:100;not null" json:"name" validate:"required"`
	Description string    `gorm:"size:255" json:"description" validate:"required"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Products    []Product `gorm:"foreignKey:CategoryID" json:"products"` // Связь с продуктами
}
