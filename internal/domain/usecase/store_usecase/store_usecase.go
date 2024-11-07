package storeUsecase

import (
	"errors"
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
func (s *StoreUseCase) CreateStore(store entities.Store, userID uint64) error {
	return s.storeRepo.Save(store, userID)
}

// GetStoreByID получает магазин по ID
func (s *StoreUseCase) GetStoreByID(id uint64) (entities.Store, error) {
	return s.storeRepo.FindByID(id)
}

// UpdateStore обновляет существующий магазин
func (s *StoreUseCase) UpdateStore(store entities.Store) error {
	return s.storeRepo.Update(store)
}

// DeleteStore удаляет магазин по ID
func (s *StoreUseCase) DeleteStore(id uint64) error {
	return s.storeRepo.Delete(id)
}

// GetAllStores получает все магазины
func (s *StoreUseCase) GetAllStores() ([]entities.Store, error) {
	return s.storeRepo.FindAll()
}

// IsUserStoreAdmin проверка является ли пользователь админом стора
func (s *StoreUseCase) IsUserStoreAdmin(storeID uint64, uid uint64) (bool, error) {
	storeExists, err := s.storeRepo.IsExist(storeID)
	if err != nil || !storeExists {
		return false, errors.New("store not found")
	}

	return s.storeRepo.IsUserStoreAdmin(storeID, uid)
}
