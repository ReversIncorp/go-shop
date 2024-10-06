package usecase

import (
	"marketplace/internal/domain/models"
)

// StoreUseCase определяет методы бизнес-логики для работы с магазинами
type StoreUseCase interface {
	CreateStore(store models.Store) error
	GetStoreByID(id string) (models.Store, error)
	UpdateStore(store models.Store) error
	DeleteStore(id string) error
	GetAllStores() ([]models.Store, error)
}
