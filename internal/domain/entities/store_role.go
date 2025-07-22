package entities

type StoreRole struct {
	StoreID uint64 `gorm:"primaryKey"`
	UserID  uint64 `gorm:"primaryKey"`
	IsOwner bool
	Store   Store `gorm:"foreignKey:StoreID"`
	User    User  `gorm:"foreignKey:UserID"`
}
