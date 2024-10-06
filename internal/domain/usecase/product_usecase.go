package usecase

import (
	"marketplace/internal/domain/models"
)

// ProductUseCase определяет методы бизнес-логики для работы с продуктами
type ProductUseCase interface {
	CreateProduct(product models.Product) error
	GetProductByID(id string) (models.Product, error)
	UpdateProduct(product models.Product) error
	DeleteProduct(id string) error
	GetProductsByStore(storeID string) ([]models.Product, error)
}
