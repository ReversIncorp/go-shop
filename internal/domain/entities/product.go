package entities

import "time"

type Product struct {
	ID          uint64  `gorm:"primaryKey"`
	Name        string  `gorm:"size:100;not null" json:"name" validate:"required"`
	Description string  `gorm:"size:500" json:"description" validate:"required"`
	Price       float64 `gorm:"not null" json:"price" validate:"required"`
	Quantity    int     `gorm:"not null" json:"quantity" validate:"required"`
	//CategoryID  uint     `gorm:"not null"`
	StoreID   uint64    `gorm:"not null" json:"store_id" validate:"required"` // Внешний ключ к таблице Store
	Store     Store     `gorm:"foreignKey:StoreID"`                           // Используем StoreID как внешний ключ
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
