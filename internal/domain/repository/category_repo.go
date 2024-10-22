package repository

import "marketplace/internal/domain/entities"

type CategoryRepository interface {
	Save(category entities.Category) error
	FindByID(id int64) (entities.Category, error)
	Update(category entities.Category) error
	Delete(id int64) error
	FindAllByStore(id int64) ([]entities.Category, error)
	IsBelongsToStore(categoryID, storeID int64) (bool, error)
	IsExist(id int64) (bool, error)
}
