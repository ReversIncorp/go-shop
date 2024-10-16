package repository

import (
	"errors"
	"gorm.io/gorm"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
	"sync"
	"time"
)

type storeRepositoryImpl struct {
	db *gorm.DB
	mu sync.Mutex
}

func NewStoreRepository(db *gorm.DB) repository.StoreRepository {
	return &storeRepositoryImpl{
		db: db,
	}
}

func (r *storeRepositoryImpl) Save(store entities.Store) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var existingStore entities.Store
	if err := r.db.Where("id = ?", store.ID).First(&existingStore).Error; err == nil {
		return errors.New("store already exists in the database")
	}

	store.CreatedAt = time.Now()
	store.UpdatedAt = store.CreatedAt

	if err := r.db.Create(&store).Error; err != nil {
		return err
	}

	return nil
}

func (r *storeRepositoryImpl) FindByID(id uint64) (entities.Store, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var store entities.Store
	if err := r.db.First(&store, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Store{}, errors.New("store not found")
		}
		return entities.Store{}, err
	}

	return store, nil
}

func (r *storeRepositoryImpl) Update(store entities.Store) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var existingStore entities.Store
	if err := r.db.First(&existingStore, store.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("store not found")
		}
		return err
	}

	store.UpdatedAt = time.Now()

	if err := r.db.Save(&store).Error; err != nil {
		return err
	}

	return nil
}

func (r *storeRepositoryImpl) Delete(id uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.db.Delete(&entities.Store{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("store not found")
		}
		return err
	}

	return nil
}

func (r *storeRepositoryImpl) FindAll() ([]entities.Store, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var stores []entities.Store
	if err := r.db.Find(&stores).Error; err != nil {
		return nil, err
	}

	return stores, nil
}
