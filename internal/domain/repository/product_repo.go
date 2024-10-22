package repository

import (
	"marketplace/internal/domain/entities"
)

type ProductRepository interface {
	Save(product entities.Product) error
	FindByID(id int64) (entities.Product, error)
	Update(product entities.Product) error
	Delete(id int64) error
	FindProductsByParams(params entities.ProductSearchParams) ([]entities.Product, error)
}
