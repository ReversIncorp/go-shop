package entities

import "time"

type Store struct {
	ID          uint64    `gorm:"primaryKey"`
	Name        string    `gorm:"size:100;not null" json:"name" validate:"required"`
	Description string    `gorm:"type:text" json:"description" validate:"required"`
	OwnerID     uint64    `gorm:"not null" json:"owner_id" validate:"required"` // Внешний ключ для связи с пользователем
	Owner       User      `gorm:"foreignKey:OwnerID" json:"owner"`              // Связь с владельцем
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Products    []Product `gorm:"foreignKey:StoreID" json:"products"` // Продукты, связанные с магазином
}
