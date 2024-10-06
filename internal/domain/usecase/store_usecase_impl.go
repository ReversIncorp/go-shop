package usecase

import (
	"marketplace/internal/domain/models"
	"marketplace/internal/domain/repository"
)

// storeUseCaseImpl реализует интерфейс StoreUseCase
type storeUseCaseImpl struct {
	storeRepo repository.StoreRepository
}

// NewStoreUseCase создает новый экземпляр StoreUseCase
func NewStoreUseCase(storeRepo repository.StoreRepository) StoreUseCase {
	return &storeUseCaseImpl{storeRepo: storeRepo}
}

// CreateStore создает новый магазин
func (s *storeUseCaseImpl) CreateStore(store models.Store) error {
	return s.storeRepo.Save(store)
}

// GetStoreByID получает магазин по ID
func (s *storeUseCaseImpl) GetStoreByID(id string) (models.Store, error) {
	return s.storeRepo.FindByID(id)
}

// UpdateStore обновляет существующий магазин
func (s *storeUseCaseImpl) UpdateStore(store models.Store) error {
	return s.storeRepo.Update(store)
}

// DeleteStore удаляет магазин по ID
func (s *storeUseCaseImpl) DeleteStore(id string) error {
	return s.storeRepo.Delete(id)
}

// GetAllStores получает все магазины
func (s *storeUseCaseImpl) GetAllStores() ([]models.Store, error) {
	return s.storeRepo.FindAll()
}
