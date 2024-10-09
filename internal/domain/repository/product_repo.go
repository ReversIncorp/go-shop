package repository

import (
	"marketplace/internal/domain/entities"
)

type ProductRepository interface {
	Save(product entities.Product) error
	FindByID(id string) (entities.Product, error)
	Update(product entities.Product) error
	Delete(id string) error
	FindAllByStore(storeID string) ([]entities.Product, error)
}
