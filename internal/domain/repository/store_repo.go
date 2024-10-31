package repository

import (
	"marketplace/internal/domain/entities"
)

type StoreRepository interface {
	Save(store entities.Store) (uint64, error)
	Update(store entities.Store) error
	Delete(id uint64) error
	IsExist(id uint64) (bool, error)

	FindByID(id uint64) (entities.Store, error)
	FindAll() ([]entities.Store, error)

	AttachCategory(storeID, categoryID uint64) error
	IsCategoryAttached(storeID, categoryID uint64) (bool, error)
	DetachCategory(storeID, categoryID uint64) error
}
