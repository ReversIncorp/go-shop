package entities

type User struct {
	ID       uint64  `gorm:"primaryKey" json:"id"`
	Name     string  `json:"name" validate:"required"`
	Email    string  `gorm:"uniqueIndex" json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=8"`
	IsSeller bool    `json:"is_seller" validate:"required"`
	Stores   []Store `gorm:"many2many:store_roles;joinForeignKey:UserID;joinReferences:StoreID" json:"stores,omitempty"`
}
