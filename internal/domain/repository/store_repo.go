package repository

import (
	"marketplace/internal/domain/entities"
)

type StoreRepository interface {
	Save(store entities.Store) error
	FindByID(id string) (entities.Store, error)
	Update(store entities.Store) error
	Delete(id string) error
	FindAll() ([]entities.Store, error)
}
