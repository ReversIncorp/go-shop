package repository

import (
	"errors"
	"fmt"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
	errorHandling "marketplace/pkg/error_handling"
	"time"

	"github.com/ztrue/tracerr"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type storeRepositoryImpl struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) repository.StoreRepository {
	return &storeRepositoryImpl{
		db: db,
	}
}

func (r *storeRepositoryImpl) IsExist(id uint64) (bool, error) {
	var store entities.Store
	err := r.db.First(&store, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, tracerr.Wrap(errors.New("store not found"))
	}
	if err != nil {
		return false, tracerr.Wrap(err)
	}
	return true, nil
}

func (r *storeRepositoryImpl) Save(store *entities.Store, uid uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		store.UpdatedAt = time.Now()
		store.CreatedAt = time.Now()
		if err := tx.Create(store).Error; err != nil {
			return tracerr.Wrap(fmt.Errorf("failed to save store: %w", err))
		}
		role := map[string]interface{}{
			"store_id": store.ID,
			"user_id":  uid,
			"is_owner": true,
		}
		if err := tx.Table("store_roles").Create(role).Error; err != nil {
			return tracerr.Wrap(fmt.Errorf("failed to add user as store admin: %w", err))
		}
		return nil
	})
}

func (r *storeRepositoryImpl) FindByID(id uint64) (entities.Store, error) {
	var store entities.Store
	err := r.db.
		First(&store, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.Store{}, errorHandling.ErrStoreNotFound
	}
	if err != nil {
		return entities.Store{}, tracerr.Wrap(err)
	}
	return store, nil
}

func (r *storeRepositoryImpl) Update(store entities.Store) error {
	store.UpdatedAt = time.Now()
	if err := r.db.Model(&store).Updates(map[string]interface{}{
		"name":        store.Name,
		"description": store.Description,
		"updated_at":  store.UpdatedAt,
	}).Error; err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func (r *storeRepositoryImpl) Delete(id uint64) error {
	if err := r.db.Delete(&entities.Store{}, id).Error; err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func (r *storeRepositoryImpl) FindAll() ([]entities.Store, error) {
	var stores []entities.Store
	if err := r.db.
		Find(&stores).Error; err != nil {
		return nil, tracerr.Wrap(err)
	}
	return stores, nil
}

func (r *storeRepositoryImpl) FindStoresByParams(params entities.StoreSearchParams) ([]entities.Store, *uint64, error) {
	db := r.db.Model(&entities.Store{})

	if params.CategoryID != nil {
		db = db.Joins("JOIN store_categories ON stores.id = store_categories.store_id").
			Where("store_categories.category_id = ?", *params.CategoryID)
	}

	if params.Name != nil {
		db = db.Where("stores.name ILIKE ?", "%"+*params.Name+"%")
	}
	if params.Cursor != nil {
		db = db.Where("stores.id > ?", *params.Cursor)
	}
	if params.Limit != nil {
		db = db.Limit(int(*params.Limit))
	}


	var stores []entities.Store
	if err := db.Order("stores.id ASC").Find(&stores).Error; err != nil {
		return nil, nil, tracerr.Wrap(fmt.Errorf("error executing query: %w", err))
	}
	var lastCursor *uint64
	if len(stores) > 0 {
		lastCursor = &stores[len(stores)-1].ID
	}
	return stores, lastCursor, nil
}

func (r *storeRepositoryImpl) AttachCategory(storeID, categoryID uint64) error {
	var store entities.Store
	if err := r.db.First(&store, storeID).Error; err != nil {
		return tracerr.Wrap(err)
	}
	var category entities.Category
	if err := r.db.First(&category, categoryID).Error; err != nil {
		return tracerr.Wrap(err)
	}
	if err := r.db.Model(&store).Association("Categories").Append(&category); err != nil {
		return tracerr.Wrap(fmt.Errorf("failed to add category to store: %w", err))
	}
	return nil
}

func (r *storeRepositoryImpl) IsCategoryAttached(storeID, categoryID uint64) (bool, error) {
	var store entities.Store
	if err := r.db.Preload("Categories", "id = ?", categoryID).First(&store, storeID).Error; err != nil {
		return false, tracerr.Wrap(err)
	}
	return len(store.Categories) > 0, nil
}

func (r *storeRepositoryImpl) DetachCategory(storeID, categoryID uint64) error {
	var store entities.Store
	if err := r.db.First(&store, storeID).Error; err != nil {
		return tracerr.Wrap(err)
	}
	var category entities.Category
	if err := r.db.First(&category, categoryID).Error; err != nil {
		return tracerr.Wrap(err)
	}
	if err := r.db.Model(&store).Association("Categories").Delete(&category); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func (r *storeRepositoryImpl) IsUserStoreAdmin(storeID, uid uint64) (bool, error) {
	var count int64
	if err := r.db.Table("store_roles").Where("store_id = ? AND user_id = ?", storeID, uid).Count(&count).Error; err != nil {
		return false, tracerr.Wrap(fmt.Errorf("failed to check if user is admin: %w", err))
	}
	return count > 0, nil
}

func (r *storeRepositoryImpl) AddUserStoreAdmin(storeID, uid uint64, owner bool) error {
	entry := map[string]interface{}{
		"store_id": storeID,
		"user_id":  uid,
		"is_owner": owner,
	}
	if err := r.db.Table("store_roles").Clauses(clause.OnConflict{DoNothing: true}).Create(entry).Error; err != nil {
		return tracerr.Wrap(fmt.Errorf("failed to add user as admin to store: %w", err))
	}
	return nil
}
