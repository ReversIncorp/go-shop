package repository

import (
	"marketplace/internal/domain/entities"
)

type ProductRepository interface {
	Save(product entities.Product, uid int64) error
	FindByID(id int64) (entities.Product, error)
	Update(product entities.Product, uid int64) error
	Delete(id int64, uid int64) error
	FindAllByStore(storeID int64) ([]entities.Product, error)
	FindAllByStoreAndCategory(storeID int64, categoryID int64) ([]entities.Product, error)
}
