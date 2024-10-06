package repository

import (
	"marketplace/internal/domain/models"
)

type StoreRepository interface {
	Save(store models.Store) error
	FindByID(id string) (models.Store, error)
	Update(store models.Store) error
	Delete(id string) error
	FindAll() ([]models.Store, error)
}
