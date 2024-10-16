package entities

import "time"

type User struct {
	ID        uint64    `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null" json:"name" validate:"required"`
	Email     string    `gorm:"size:255;unique;not null" json:"email" validate:"required,email"`
	Password  string    `gorm:"size:255;not null" json:"password" validate:"required,min=8"`
	IsSeller  bool      `gorm:"default:false" json:"is_seller"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Stores    []Store   `gorm:"foreignKey:OwnerID" json:"stores"` // Связь с магазинами
}
