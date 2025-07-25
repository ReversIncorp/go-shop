package entities

type StoreCategory struct {
	StoreID    uint64   `gorm:"primaryKey"`
	CategoryID uint64   `gorm:"primaryKey"`
	Store      Store    `gorm:"foreignKey:StoreID"`
	Category   Category `gorm:"foreignKey:CategoryID"`
}
