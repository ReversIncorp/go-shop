package repository

import (
	"errors"
	"marketplace/internal/domain/models"
	repository2 "marketplace/internal/domain/repository"
	"sync"
	"time"
)

type inMemoryStoreRepository struct {
	stores map[string]models.Store
	mu     sync.Mutex
}

func NewStoreRepository() repository2.StoreRepository {
	return &inMemoryStoreRepository{
		stores: make(map[string]models.Store),
	}
}

func (r *inMemoryStoreRepository) Save(store models.Store) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Если магазин уже существует, вернем ошибку
	if _, exists := r.stores[store.ID]; exists {
		return errors.New("store already exists")
	}

	// Устанавливаем время создания и обновления
	store.CreatedAt = time.Now()
	store.UpdatedAt = store.CreatedAt

	r.stores[store.ID] = store
	return nil
}

func (r *inMemoryStoreRepository) FindByID(id string) (models.Store, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	store, exists := r.stores[id]
	if !exists {
		return models.Store{}, errors.New("store not found")
	}

	return store, nil
}

func (r *inMemoryStoreRepository) Update(store models.Store) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.stores[store.ID]
	if !exists {
		return errors.New("store not found")
	}

	// Обновляем время изменения
	store.UpdatedAt = time.Now()

	r.stores[store.ID] = store
	return nil
}

func (r *inMemoryStoreRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.stores[id]
	if !exists {
		return errors.New("store not found")
	}

	delete(r.stores, id)
	return nil
}

func (r *inMemoryStoreRepository) FindAll() ([]models.Store, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var stores []models.Store
	for _, store := range r.stores {
		stores = append(stores, store)
	}

	return stores, nil
}
