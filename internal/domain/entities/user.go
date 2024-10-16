package entities

import "time"

type User struct {
	ID        uint64    `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null" json:"name" validate:"required"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password  string    `gorm:"not null" json:"password" validate:"required,min=8"`
	IsSeller  bool      `gorm:"default:false"` // Продавец или покупатель
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Stores    []Store   `gorm:"foreignKey:OwnerID"` // Один пользователь может владеть несколькими магазинами
}
