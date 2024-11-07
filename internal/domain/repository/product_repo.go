package repository

import (
	"marketplace/internal/domain/entities"
)

type ProductRepository interface {
	Save(product entities.Product) error
	Update(product entities.Product) error
	Delete(id uint64) error

	FindByID(id uint64) (entities.Product, error)
	FindProductsByParams(params entities.ProductSearchParams) ([]entities.Product, error)

	IsProductBelongsToStore(productID, storeID uint64) (bool, error)
}
