package repository

import (
	"marketplace/internal/domain/entities"
)

type ProductRepository interface {
	Save(product entities.Product) error
	FindByID(id int64) (entities.Product, error)
	Update(product entities.Product) error
	Delete(id int64) error
	FindAllByStore(storeID int64) ([]entities.Product, error)
	FindAllByStoreAndCategory(storeID int64, categoryID int64) ([]entities.Product, error)
}
