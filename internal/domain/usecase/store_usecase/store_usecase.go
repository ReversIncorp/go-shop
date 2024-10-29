package storeUsecase

import (
	"errors"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
)

// StoreUseCase реализует интерфейс StoreUseCase
type StoreUseCase struct {
	storeRepo repository.StoreRepository
	userRepo  repository.UserRepository
}

// NewStoreUseCase создает новый экземпляр StoreUseCase
func NewStoreUseCase(storeRepo repository.StoreRepository, userRepo repository.UserRepository) *StoreUseCase {
	return &StoreUseCase{storeRepo: storeRepo, userRepo: userRepo}
}

// CreateStore создает новый магазин
func (s *StoreUseCase) CreateStore(store entities.Store, userID uint64) error {
	storeID, err := s.storeRepo.Save(store)
	if err != nil {
		return err
	}

	return s.userRepo.AddOwningStore(userID, storeID)
}

// GetStoreByID получает магазин по ID
func (s *StoreUseCase) GetStoreByID(id uint64) (entities.Store, error) {
	return s.storeRepo.FindByID(id)
}

// UpdateStore обновляет существующий магазин
func (s *StoreUseCase) UpdateStore(store entities.Store, uid uint64) error {
	storeExists, err := s.storeRepo.IsExist(store.ID)
	if err != nil || !storeExists {
		return errors.New("store not found")
	}

	isOwner, err := s.userRepo.IsOwnsStore(uid, store.ID)
	if err != nil || !isOwner {
		return errors.New("user does not owning this store")
	}

	return s.storeRepo.Update(store)
}

// DeleteStore удаляет магазин по ID
func (s *StoreUseCase) DeleteStore(id uint64, uid uint64) error {
	storeExists, err := s.storeRepo.IsExist(id)
	if err != nil || !storeExists {
		return errors.New("store not found")
	}

	isOwner, err := s.userRepo.IsOwnsStore(uid, id)
	if err != nil || !isOwner {
		return errors.New("user does not owning this store")
	}

	return s.storeRepo.Delete(id)
}

// GetAllStores получает все магазины
func (s *StoreUseCase) GetAllStores() ([]entities.Store, error) {
	return s.storeRepo.FindAll()
}
