package usecase

import (
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
)

// StoreUseCase реализует интерфейс StoreUseCase
type StoreUseCase struct {
	storeRepo repository.StoreRepository
}

// NewStoreUseCase создает новый экземпляр StoreUseCase
func NewStoreUseCase(storeRepo repository.StoreRepository) *StoreUseCase {
	return &StoreUseCase{storeRepo: storeRepo}
}

// CreateStore создает новый магазин
func (s *StoreUseCase) CreateStore(store entities.Store) error {
	return s.storeRepo.Save(store)
}

// GetStoreByID получает магазин по ID
func (s *StoreUseCase) GetStoreByID(id string) (entities.Store, error) {
	return s.storeRepo.FindByID(id)
}

// UpdateStore обновляет существующий магазин
func (s *StoreUseCase) UpdateStore(store entities.Store) error {
	return s.storeRepo.Update(store)
}

// DeleteStore удаляет магазин по ID
func (s *StoreUseCase) DeleteStore(id string) error {
	return s.storeRepo.Delete(id)
}

// GetAllStores получает все магазины
func (s *StoreUseCase) GetAllStores() ([]entities.Store, error) {
	return s.storeRepo.FindAll()
}
