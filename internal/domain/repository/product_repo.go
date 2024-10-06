package repository

import (
	"marketplace/internal/domain/models"
)

type ProductRepository interface {
	Save(product models.Product) error
	FindByID(id string) (models.Product, error)
	Update(product models.Product) error
	Delete(id string) error
	FindAllByStore(storeID string) ([]models.Product, error)
}
