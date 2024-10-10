package repository

import (
	"marketplace/internal/domain/entities"
)

type StoreRepository interface {
	Save(store entities.Store) error
	FindByID(id uint64) (entities.Store, error)
	Update(store entities.Store) error
	Delete(id uint64) error
	FindAll() ([]entities.Store, error)
}
