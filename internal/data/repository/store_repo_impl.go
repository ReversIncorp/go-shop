package repository

import (
	"errors"
	"marketplace/internal/domain/entities"
	repository2 "marketplace/internal/domain/repository"
	"sync"
	"time"
)

type inMemoryStoreRepository struct {
	stores map[string]entities.Store
	mu     sync.Mutex
}

func NewStoreRepository() repository2.StoreRepository {
	return &inMemoryStoreRepository{
		stores: make(map[string]entities.Store),
	}
}

func (r *inMemoryStoreRepository) Save(store entities.Store) error {
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

func (r *inMemoryStoreRepository) FindByID(id string) (entities.Store, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	store, exists := r.stores[id]
	if !exists {
		return entities.Store{}, errors.New("store not found")
	}

	return store, nil
}

func (r *inMemoryStoreRepository) Update(store entities.Store) error {
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

func (r *inMemoryStoreRepository) FindAll() ([]entities.Store, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var stores []entities.Store
	for _, store := range r.stores {
		stores = append(stores, store)
	}

	return stores, nil
}
