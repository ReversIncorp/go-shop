package entities

import "time"

type Store struct {
	ID          uint64     `gorm:"primaryKey" json:"id"`
	Name        string     `json:"name" validate:"required"`
	Description string     `json:"description" validate:"required"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Products    []Product  `gorm:"foreignKey:StoreID" json:"products,omitempty"`
	Categories  []Category `gorm:"many2many:store_categories" json:"categories,omitempty"`
	Users       []User     `gorm:"many2many:store_roles;joinForeignKey:StoreID;joinReferences:UserID" json:"users,omitempty"`
}
