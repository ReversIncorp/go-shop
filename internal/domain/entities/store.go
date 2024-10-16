package entities

import "time"

type Store struct {
	ID          uint64 `json:"id" json:"id"`
	Name        string `json:"name" json:"name" validate:"required"`
	Description string `json:"description" json:"description" validate:"required"`
	OwnerID     uint   `gorm:"not null"  json:"owner_id" validate:"required"`
	//Owner       User      `gorm:"foreignKey:OwnerID"`
	Products  []Product `gorm:"foreignKey:StoreID"` // Используем StoreID для связи с продуктами
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
