package repository

import (
	"marketplace/internal/domain/entities"
)

type ProductRepository interface {
	Save(product entities.Product) error
	FindByID(id uint64) (entities.Product, error)
	Update(product entities.Product) error
	Delete(id uint64) error
	FindAllByStore(storeID uint64) ([]entities.Product, error)
}
