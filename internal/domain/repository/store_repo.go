package repository

import (
	"marketplace/internal/domain/entities"
)

type StoreRepository interface {
	Save(store entities.Store) (int64, error)
	FindByID(id int64) (entities.Store, error)
	Update(store entities.Store) error
	Delete(id int64) error
	FindAll() ([]entities.Store, error)
	IsExist(id int64) (bool, error)
}
