package repository

import (
	"marketplace/internal/domain/entities"
)

type StoreRepository interface {
	Save(store entities.Store, uid uint64) error
	Update(store entities.Store) error
	Delete(id uint64) error
	IsExist(id uint64) (bool, error)

	FindByID(id uint64) (entities.Store, error)
	FindAll() ([]entities.Store, error)

	IsUserStoreAdmin(storeID, uid uint64) (bool, error)
	AddUserStoreAdmin(storeID, uid uint64, owner bool) error
}
